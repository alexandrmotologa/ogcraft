import React, { useState } from 'react';
import { Template } from '../lib/types';
import { ExternalLink, Heart, MessageCircle, Repeat2, Share, Bookmark, ThumbsUp, MessageSquare, Send } from 'lucide-react';

interface SocialFeedSimulatorProps {
  template: Template;
  previewUrl: string | null;
  params: Record<string, string>;
}

export const SocialFeedSimulator: React.FC<SocialFeedSimulatorProps> = ({
  template,
  previewUrl,
  params,
}) => {
  const [platform, setPlatform] = useState<'twitter' | 'linkedin' | 'discord' | 'slack'>('twitter');

  const title = params.title || template.elements.find(e => e.field === 'title')?.default || 'Building Modern Systems';
  const author = params.author || template.elements.find(e => e.field === 'author')?.default || 'Alexandr Motologa';
  const tag = params.tag || template.elements.find(e => e.field === 'tag')?.default || 'Engineering';

  return (
    <div className="flex-1 flex flex-col h-full bg-[#070a10] overflow-y-auto p-6">
      {/* Platform Switcher */}
      <div className="flex items-center justify-between mb-6 pb-4 border-b border-[#1a2333]">
        <div>
          <h2 className="text-sm font-bold text-slate-200">Social Feed Simulator</h2>
          <p className="text-xs text-slate-500">Preview how social crawlers and feeds unfurl this card</p>
        </div>

        <div className="flex items-center bg-[#0c121e] border border-slate-800 rounded-lg p-1 gap-1">
          <button
            onClick={() => setPlatform('twitter')}
            className={`px-3 py-1.5 text-xs font-semibold rounded-md transition ${
              platform === 'twitter'
                ? 'bg-sky-500 text-slate-950 shadow'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Twitter / X
          </button>
          <button
            onClick={() => setPlatform('linkedin')}
            className={`px-3 py-1.5 text-xs font-semibold rounded-md transition ${
              platform === 'linkedin'
                ? 'bg-[#0a66c2] text-white shadow'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            LinkedIn
          </button>
          <button
            onClick={() => setPlatform('discord')}
            className={`px-3 py-1.5 text-xs font-semibold rounded-md transition ${
              platform === 'discord'
                ? 'bg-[#5865f2] text-white shadow'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Discord
          </button>
          <button
            onClick={() => setPlatform('slack')}
            className={`px-3 py-1.5 text-xs font-semibold rounded-md transition ${
              platform === 'slack'
                ? 'bg-[#4a154b] text-white shadow'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Slack
          </button>
        </div>
      </div>

      {/* Mockup Preview Area */}
      <div className="flex-1 flex items-center justify-center">
        {/* Twitter / X Mockup */}
        {platform === 'twitter' && (
          <div className="w-full max-w-[560px] bg-black border border-slate-800 rounded-2xl p-4 shadow-2xl">
            {/* Tweet Author Header */}
            <div className="flex items-center gap-3 mb-3">
              <div className="w-10 h-10 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center font-bold text-sky-400 text-sm">
                {author.slice(0, 2).toUpperCase()}
              </div>
              <div className="flex flex-col">
                <div className="flex items-center gap-1.5">
                  <span className="font-bold text-sm text-slate-100">{author}</span>
                  <span className="text-slate-500 text-xs">@mtlglabs</span>
                  <span className="text-slate-500 text-xs">· 2h</span>
                </div>
                <span className="text-xs text-slate-400">Just shipped a comprehensive deep dive on {tag}. Check it out:</span>
              </div>
            </div>

            {/* Twitter Large Summary Card */}
            <div className="rounded-2xl border border-slate-800 overflow-hidden bg-[#0a0e17] transition hover:border-slate-700 cursor-pointer">
              {previewUrl ? (
                <img src={previewUrl} alt="Twitter Preview" className="w-full aspect-[1.904/1] object-cover" />
              ) : (
                <div className="w-full aspect-[1.904/1] bg-slate-900 animate-pulse" />
              )}
              <div className="p-3">
                <span className="text-[11px] font-mono text-slate-500 block truncate">mtlglabs.space</span>
                <span className="text-sm font-semibold text-slate-200 line-clamp-1 mt-0.5">{title}</span>
                <span className="text-xs text-slate-400 line-clamp-1 mt-0.5">Real-time vector rendering in Go with zero browser dependencies.</span>
              </div>
            </div>

            {/* Tweet Actions */}
            <div className="flex items-center justify-between text-slate-500 text-xs px-2 pt-3 mt-1">
              <span className="flex items-center gap-1.5 hover:text-sky-400"><MessageCircle size={15} /> 24</span>
              <span className="flex items-center gap-1.5 hover:text-emerald-400"><Repeat2 size={16} /> 48</span>
              <span className="flex items-center gap-1.5 hover:text-rose-400"><Heart size={15} /> 312</span>
              <span className="flex items-center gap-1.5 hover:text-sky-400"><Bookmark size={15} /></span>
              <span className="flex items-center gap-1.5 hover:text-sky-400"><Share size={15} /></span>
            </div>
          </div>
        )}

        {/* LinkedIn Mockup */}
        {platform === 'linkedin' && (
          <div className="w-full max-w-[550px] bg-[#1b1f23] border border-slate-700/60 rounded-xl p-4 shadow-2xl">
            <div className="flex items-center gap-3 mb-3">
              <div className="w-11 h-11 rounded-full bg-slate-800 border border-slate-700 flex items-center justify-center font-bold text-sky-400">
                {author.slice(0, 2).toUpperCase()}
              </div>
              <div>
                <div className="font-bold text-sm text-slate-100">{author}</div>
                <div className="text-[11px] text-slate-400">Software Architect • 1st</div>
                <div className="text-[10px] text-slate-500">2h • Edited • 🌐</div>
              </div>
            </div>

            <p className="text-xs text-slate-300 leading-relaxed mb-3">
              Thrilled to share our latest milestone with {title}. Native SIMD vector execution delivers sub-15ms rendering directly from URL query parameters.
            </p>

            {/* LinkedIn Card */}
            <div className="rounded-lg border border-slate-700 overflow-hidden bg-[#24292f] cursor-pointer">
              {previewUrl ? (
                <img src={previewUrl} alt="LinkedIn Preview" className="w-full aspect-[1.904/1] object-cover" />
              ) : (
                <div className="w-full aspect-[1.904/1] bg-slate-900" />
              )}
              <div className="p-3">
                <span className="text-sm font-bold text-slate-100 line-clamp-1">{title}</span>
                <span className="text-[11px] text-slate-400 font-mono block mt-0.5">mtlglabs.space • 4 min read</span>
              </div>
            </div>

            {/* LinkedIn Footer Reactions */}
            <div className="flex items-center justify-between text-slate-400 text-xs pt-3 mt-2 border-t border-slate-700/60">
              <span className="flex items-center gap-1.5 hover:text-sky-400"><ThumbsUp size={15} /> Like</span>
              <span className="flex items-center gap-1.5 hover:text-sky-400"><MessageSquare size={15} /> Comment</span>
              <span className="flex items-center gap-1.5 hover:text-sky-400"><Repeat2 size={15} /> Repost</span>
              <span className="flex items-center gap-1.5 hover:text-sky-400"><Send size={15} /> Send</span>
            </div>
          </div>
        )}

        {/* Discord Embed Mockup */}
        {platform === 'discord' && (
          <div className="w-full max-w-[500px] bg-[#313338] rounded-lg p-4 shadow-2xl font-sans">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-8 h-8 rounded-full bg-[#5865f2] flex items-center justify-center font-bold text-white text-xs">
                BOT
              </div>
              <span className="text-sm font-semibold text-slate-200">SocialBot</span>
              <span className="text-[10px] px-1 py-0.2 rounded bg-[#5865f2] text-white uppercase font-bold">APP</span>
              <span className="text-[10px] text-slate-400">Today at 4:18 PM</span>
            </div>

            {/* Discord Embed Box */}
            <div className="bg-[#2b2d31] border-l-4 border-[#38bdf8] rounded-r p-3.5 space-y-2">
              <span className="text-[11px] text-slate-400 font-medium">mtlglabs.space</span>
              <a href="#" className="text-sm font-bold text-[#00a8fc] hover:underline block">
                {title}
              </a>
              <p className="text-xs text-slate-300">
                Self-hosted dynamic OpenGraph image and social preview generator with two-tier caching.
              </p>
              {previewUrl && (
                <img src={previewUrl} alt="Discord Preview" className="w-full rounded-md object-cover max-h-72 mt-2" />
              )}
            </div>
          </div>
        )}

        {/* Slack Unfurl Mockup */}
        {platform === 'slack' && (
          <div className="w-full max-w-[520px] bg-[#1a1d21] rounded-lg p-4 shadow-2xl">
            <div className="flex items-center gap-2 mb-2">
              <span className="font-bold text-sm text-slate-200">alexandr</span>
              <span className="text-[11px] text-slate-500">4:20 PM</span>
            </div>
            <p className="text-xs text-sky-400 mb-2 underline cursor-pointer">
              https://mtlglabs.space/blog/{template.id}
            </p>

            {/* Slack Attachment Card */}
            <div className="border-l-4 border-[#38bdf8] pl-3 py-1 space-y-1.5">
              <div className="flex items-center gap-1.5 text-xs font-bold text-slate-200">
                <span>{title}</span>
                <ExternalLink size={12} className="text-slate-500" />
              </div>
              <p className="text-[11px] text-slate-400">
                Sub-15ms card rendering in Go with Inter and JetBrains Mono typography.
              </p>
              {previewUrl && (
                <img src={previewUrl} alt="Slack Preview" className="w-full rounded-md max-w-[440px] mt-2 border border-slate-800" />
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
