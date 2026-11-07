# ogcraft Integration Guides

Guides for integrating `ogcraft` dynamic preview cards into modern web frameworks and static site generators.

---

## 1. Next.js 15 App Router

Use `generateMetadata` in your dynamic page route (e.g. `app/blog/[slug]/page.tsx`):

```typescript
import type { Metadata } from "next";

interface Props {
  params: Promise<{ slug: string }>;
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const post = await getPost(slug);

  const ogUrl = new URL("https://og.yourdomain.com/v1/cards/dark_tech.png");
  ogUrl.searchParams.set("title", post.title);
  ogUrl.searchParams.set("author", post.author.name);
  ogUrl.searchParams.set("tag", post.category);
  ogUrl.searchParams.set("avatar", post.author.avatarUrl);

  return {
    title: post.title,
    description: post.excerpt,
    openGraph: {
      title: post.title,
      description: post.excerpt,
      images: [
        {
          url: ogUrl.toString(),
          width: 1200,
          height: 630,
          alt: post.title,
        },
      ],
    },
    twitter: {
      card: "summary_large_image",
      title: post.title,
      images: [ogUrl.toString()],
    },
  };
}
```

---

## 2. Astro

In your blog post template (`src/layouts/BlogPost.astro`):

```astro
---
const { title, description, author, tag, avatar } = Astro.props;

const ogParams = new URLSearchParams({
  title,
  author,
  tag,
  avatar: avatar || "https://yourdomain.com/default-avatar.png",
});

const ogImage = `https://og.yourdomain.com/v1/cards/minimal_light.png?${ogParams.toString()}`;
---

<head>
  <meta property="og:title" content={title} />
  <meta property="og:description" content={description} />
  <meta property="og:image" content={ogImage} />
  <meta property="og:image:width" content="1200" />
  <meta property="og:image:height" content="630" />
  
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:image" content={ogImage} />
</head>
```

---

## 3. Hugo

In `layouts/partials/head.html`:

```html
{{ $ogBase := "https://og.yourdomain.com/v1/cards/dark_tech.png" }}
{{ $title := .Title }}
{{ $tag := (index .Params.tags 0) | default "Article" }}
{{ $author := .Site.Params.author | default "Alexandr Motologa" }}

{{ $ogImage := printf "%s?title=%s&tag=%s&author=%s" $ogBase (querify "title" $title) (querify "tag" $tag) (querify "author" $author) }}

<meta property="og:image" content="{{ $ogImage }}">
<meta property="og:image:width" content="1200">
<meta property="og:image:height" content="630">
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:image" content="{{ $ogImage }}">
```

---

## 4. Plain HTML / Static Markdown

Add standard OpenGraph meta tags:

```html
<meta property="og:image" content="https://og.yourdomain.com/v1/cards/dark_tech.png?title=High+Performance+Go+Services&author=Alexandr&tag=Go">
<meta property="og:image:width" content="1200">
<meta property="og:image:height" content="630">
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:image" content="https://og.yourdomain.com/v1/cards/dark_tech.png?title=High+Performance+Go+Services&author=Alexandr&tag=Go">
```
