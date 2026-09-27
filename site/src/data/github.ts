// Build-time GitHub stats. Fetched once during `astro build` and inlined into the
// static HTML — no client JS, no runtime rate limits. The daily deploy cron
// (.github/workflows/deploy.yaml) re-runs the build so the number stays fresh.
// Falls back gracefully so an offline / rate-limited build never fails.
const REPO = 'labstack/echo';
const FALLBACK_STARS = 32400;

async function fetchStars(): Promise<number> {
  try {
    const headers: Record<string, string> = {
      Accept: 'application/vnd.github+json',
      'User-Agent': 'echo-docs-build',
    };
    // Optional: set GITHUB_TOKEN in CI to lift the 60 req/hr unauthenticated limit.
    const token = process.env.GITHUB_TOKEN;
    if (token) headers.Authorization = `Bearer ${token}`;

    const res = await fetch(`https://api.github.com/repos/${REPO}`, { headers, signal: AbortSignal.timeout(5000) });
    if (!res.ok) throw new Error(`GitHub API responded ${res.status}`);
    const data = await res.json();
    return typeof data.stargazers_count === 'number' ? data.stargazers_count : FALLBACK_STARS;
  } catch (e) {
    console.warn(`[github] star fetch failed; using fallback ${FALLBACK_STARS}:`, e);
    return FALLBACK_STARS;
  }
}

export const stars = await fetchStars();

/** e.g. 32412 -> "32.4k" */
export const starsLabel = stars >= 1000 ? `${(stars / 1000).toFixed(1)}k` : String(stars);
