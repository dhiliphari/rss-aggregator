const PAGE_SIZE = 30;
const AUTO_REFRESH_MS = 60_000;

function stripHtml(html) {
  return new DOMParser().parseFromString(html || "", "text/html").body.textContent.trim();
}

document.addEventListener("alpine:init", () => {
  Alpine.data("app", () => ({
    apiKey: localStorage.getItem("apiKey"),
    user: null,
    booting: !!localStorage.getItem("apiKey"),

    feeds: [],
    follows: [],
    posts: [],
    postCount: PAGE_SIZE,

    authTab: "login",
    feedTab: "following",
    feedQuery: "",
    postQuery: "",
    postFeed: null,
    showAddFeed: false,
    newApiKey: null,

    forms: { name: "", apiKey: "", feedName: "", feedUrl: "" },
    busy: { auth: false, feed: false, posts: false, follow: {} },
    toasts: [],

    init() {
      if (this.apiKey) {
        this.login(this.apiKey)
          .catch((err) => this.toast(err.message, "error"))
          .finally(() => (this.booting = false));
      }
      setInterval(() => {
        if (this.user && !document.hidden) this.loadPosts().catch(() => {});
      }, AUTO_REFRESH_MS);
    },

    // --- derived state ---

    get followByFeed() {
      return new Map(this.follows.map((f) => [f.feed_id, f.id]));
    },

    get followedFeeds() {
      return this.feeds.filter((f) => this.followByFeed.has(f.id));
    },

    get visibleFeeds() {
      const list = this.feedTab === "following" ? this.followedFeeds : this.feeds;
      const q = this.feedQuery.trim().toLowerCase();
      if (!q) return list;
      return list.filter((f) => f.name.toLowerCase().includes(q) || f.url.toLowerCase().includes(q));
    },

    get visiblePosts() {
      const q = this.postQuery.trim().toLowerCase();
      return this.posts.filter((p) =>
        (!this.postFeed || p.feed_id === this.postFeed) &&
        (!q || p.title.toLowerCase().includes(q) || p.summary.toLowerCase().includes(q)),
      );
    },

    isFollowing(feedId) {
      return this.followByFeed.has(feedId);
    },

    feedName(feedId) {
      const feed = this.feeds.find((f) => f.id === feedId);
      return feed ? feed.name : "Unknown feed";
    },

    togglePostFeed(feedId) {
      this.postFeed = this.postFeed === feedId ? null : feedId;
    },

    // --- api ---

    async api(method, path, body) {
      const headers = {};
      if (this.apiKey) headers["Authorization"] = `ApiKey ${this.apiKey}`;
      if (body) headers["Content-Type"] = "application/json";

      const res = await fetch(path, {
        method,
        headers,
        body: body ? JSON.stringify(body) : undefined,
      });

      if (res.status === 204) return null;

      const text = await res.text();
      const data = text ? JSON.parse(text) : null;

      if (!res.ok) {
        let msg = data && data.error ? data.error : res.statusText;
        if (typeof msg === "object") {
          msg = Object.entries(msg).map(([k, v]) => `${k}: ${v}`).join(", ");
        }
        throw new Error(msg);
      }
      return data;
    },

    // --- auth ---

    async login(apiKey) {
      this.apiKey = apiKey;
      try {
        this.user = await this.api("GET", "/me");
      } catch (err) {
        this.logout();
        throw err;
      }
      localStorage.setItem("apiKey", apiKey);
      await this.refresh();
    },

    logout() {
      this.apiKey = null;
      this.user = null;
      this.feeds = [];
      this.follows = [];
      this.posts = [];
      this.postFeed = null;
      localStorage.removeItem("apiKey");
    },

    async submitLogin() {
      this.busy.auth = true;
      try {
        await this.login(this.forms.apiKey);
        this.forms.apiKey = "";
      } catch (err) {
        this.toast(err.message, "error");
      } finally {
        this.busy.auth = false;
      }
    },

    async submitSignup() {
      this.busy.auth = true;
      try {
        const user = await this.api("POST", "/users", { name: this.forms.name });
        await this.login(user.api_key);
        this.newApiKey = user.api_key;
        this.forms.name = "";
      } catch (err) {
        this.toast(err.message, "error");
      } finally {
        this.busy.auth = false;
      }
    },

    // --- data ---

    async refresh() {
      try {
        await Promise.all([this.loadFeeds(), this.loadPosts()]);
      } catch (err) {
        this.toast(err.message, "error");
      }
    },

    async loadFeeds() {
      const [feeds, follows] = await Promise.all([
        this.api("GET", "/feeds"),
        this.api("GET", "/me/feeds/following"),
      ]);
      this.feeds = feeds || [];
      this.follows = follows || [];
    },

    async loadPosts() {
      this.busy.posts = true;
      try {
        const posts = await this.api("GET", `/me/posts?count=${this.postCount}`);
        this.posts = (posts || []).map((p) => {
          const text = stripHtml(p.description);
          return { ...p, summary: text.length > 240 ? text.slice(0, 240) + "…" : text };
        });
      } finally {
        this.busy.posts = false;
      }
    },

    async loadMore() {
      this.postCount += PAGE_SIZE;
      await this.loadPosts().catch((err) => this.toast(err.message, "error"));
    },

    async toggleFollow(feed) {
      this.busy.follow[feed.id] = true;
      try {
        const followId = this.followByFeed.get(feed.id);
        if (followId) {
          await this.api("DELETE", `/me/feeds/following/${followId}`);
          if (this.postFeed === feed.id) this.postFeed = null;
          this.toast(`Unfollowed ${feed.name}`);
        } else {
          await this.api("POST", "/me/feeds/following", { feed_id: feed.id });
          this.toast(`Following ${feed.name}`);
        }
        await this.refresh();
      } catch (err) {
        this.toast(err.message, "error");
      } finally {
        delete this.busy.follow[feed.id];
      }
    },

    async submitFeed() {
      this.busy.feed = true;
      try {
        const feed = await this.api("POST", "/me/feeds", {
          name: this.forms.feedName,
          url: this.forms.feedUrl,
        });
        await this.api("POST", "/me/feeds/following", { feed_id: feed.id });
        this.forms.feedName = "";
        this.forms.feedUrl = "";
        this.showAddFeed = false;
        this.feedTab = "following";
        this.toast(`Added ${feed.name}. Posts will appear within a minute.`);
        await this.refresh();
      } catch (err) {
        this.toast(err.message, "error");
      } finally {
        this.busy.feed = false;
      }
    },

    // --- ui helpers ---

    toast(message, type = "info") {
      const id = Date.now() + Math.random();
      this.toasts.push({ id, message, type });
      setTimeout(() => this.dismiss(id), type === "error" ? 6000 : 3000);
    },

    dismiss(id) {
      this.toasts = this.toasts.filter((t) => t.id !== id);
    },

    async copy(text, message) {
      try {
        await navigator.clipboard.writeText(text);
        this.toast(message);
      } catch {
        this.toast("Could not copy to clipboard", "error");
      }
    },

    hostname(url) {
      try {
        return new URL(url).hostname;
      } catch {
        return url;
      }
    },

    timeAgo(date) {
      const seconds = Math.round((new Date(date) - Date.now()) / 1000);
      const units = [
        ["year", 31536000], ["month", 2592000], ["week", 604800],
        ["day", 86400], ["hour", 3600], ["minute", 60],
      ];
      const rtf = new Intl.RelativeTimeFormat(undefined, { numeric: "auto" });
      for (const [unit, secs] of units) {
        if (Math.abs(seconds) >= secs) return rtf.format(Math.round(seconds / secs), unit);
      }
      return "just now";
    },
  }));
});
