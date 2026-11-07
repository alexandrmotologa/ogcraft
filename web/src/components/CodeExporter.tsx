import React, { useState } from 'react';
import { Template } from '../lib/types';
import { buildCardUrl } from '../lib/api';
import { Code, Copy, Check, ExternalLink } from 'lucide-react';

interface CodeExporterProps {
  template: Template;
  params: Record<string, string>;
}

export const CodeExporter: React.FC<CodeExporterProps> = ({ template, params }) => {
  const [tab, setTab] = useState<'url' | 'nextjs' | 'astro' | 'html' | 'markdown'>('url');
  const [copied, setCopied] = useState<boolean>(false);

  const cardUrl = buildCardUrl(template.id, params);

  const getSnippet = () => {
    switch (tab) {
      case 'url':
        return cardUrl;

      case 'nextjs':
        return `// app/blog/[slug]/page.tsx
export async function generateMetadata({ params }): Promise<Metadata> {
  const post = await getPost(params.slug);
  const ogUrl = new URL('${window.location.origin}/v1/cards/${template.id}.png');
  ogUrl.searchParams.set('title', post.title);
  ogUrl.searchParams.set('author', post.author);
  ogUrl.searchParams.set('tag', post.category);

  return {
    title: post.title,
    openGraph: {
      images: [{ url: ogUrl.toString(), width: 1200, height: 630 }],
    },
    twitter: {
      card: 'summary_large_image',
      images: [ogUrl.toString()],
    },
  };
}`;

      case 'astro':
        return `---
// src/layouts/BlogPost.astro
const { title, author, tag } = Astro.props;
const ogParams = new URLSearchParams({ title, author, tag });
const ogImage = \`${window.location.origin}/v1/cards/${template.id}.png?\${ogParams.toString()}\`;
---
<head>
  <meta property="og:image" content={ogImage} />
  <meta property="og:image:width" content="1200" />
  <meta property="og:image:height" content="630" />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:image" content={ogImage} />
</head>`;

      case 'html':
        return `<meta property="og:image" content="${cardUrl}" />
<meta property="og:image:width" content="1200" />
<meta property="og:image:height" content="630" />
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:image" content="${cardUrl}" />`;

      case 'markdown':
        return `[![Social Card](${cardUrl})](${cardUrl})`;
    }
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(getSnippet());
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="bg-[#0c121e] border-t border-[#1a2333] p-4 flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Code size={15} className="text-sky-400" />
          <span className="text-xs font-bold uppercase tracking-wider text-slate-300">
            Export Snippets
          </span>
        </div>

        {/* Tabs */}
        <div className="flex items-center bg-slate-900 border border-slate-800 rounded-lg p-0.5">
          {(['url', 'nextjs', 'astro', 'html', 'markdown'] as const).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={`px-2.5 py-1 text-[11px] font-mono rounded transition uppercase ${
                tab === t
                  ? 'bg-sky-500/20 text-sky-300 font-semibold'
                  : 'text-slate-400 hover:text-slate-200'
              }`}
            >
              {t}
            </button>
          ))}
        </div>

        <div className="flex items-center gap-2">
          <a
            href={cardUrl}
            target="_blank"
            rel="noreferrer"
            className="flex items-center gap-1 text-[11px] text-slate-400 hover:text-sky-300 transition"
          >
            <ExternalLink size={12} />
            Test in New Tab
          </a>
          <button
            onClick={handleCopy}
            className="flex items-center gap-1.5 px-3 py-1 text-xs font-medium bg-slate-800 hover:bg-slate-700 text-slate-200 rounded transition"
          >
            {copied ? <Check size={12} className="text-emerald-400" /> : <Copy size={12} />}
            {copied ? 'Copied!' : 'Copy'}
          </button>
        </div>
      </div>

      <pre className="bg-[#070a10] border border-slate-800/80 rounded-lg p-3 text-[11px] font-mono text-slate-300 overflow-x-auto selection:bg-sky-500/30">
        <code>{getSnippet()}</code>
      </pre>
    </div>
  );
};
