import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: "6022 Agentic Team protocol",
  tagline: "A protocol for building agentic teams",
  favicon: "img/logo-192x192.png",

  // Future flags, see https://docusaurus.io/docs/api/docusaurus-config#future
  future: {
    v4: true, // Improve compatibility with the upcoming Docusaurus v4
  },

  // Set the production url of your site here
  url: "https://your-docusaurus-site.example.com",
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: "/",

  // Enable Mermaid diagrams in Markdown/MDX
  markdown: {
    mermaid: true,
    hooks: {
      onBrokenMarkdownLinks: "throw",
    },
  },

  themes: ["@docusaurus/theme-mermaid"],

  onBrokenLinks: "throw",

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      {
        docs: {
          sidebarPath: "./sidebars.ts",
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            "https://github.com/facebook/docusaurus/tree/main/packages/create-docusaurus/templates/shared/",
        },
        theme: {
          customCss: "./src/css/custom.css",
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Mermaid theme configuration (optional)
    mermaid: {
      theme: { light: "neutral", dark: "dark" },
    },
    // Replace with your project's social card
    image: "img/logo.png",
    navbar: {
      title: "6022 Agentic Team protocol",
      logo: {
        alt: "6022",
        src: "img/logo.png",
      },
      items: [
        {
          type: "docSidebar",
          sidebarId: "usersSidebar",
          position: "left",
          label: "User Guide",
        },
        {
          type: "docSidebar",
          sidebarId: "developersSidebar",
          position: "left",
          label: "Developers",
        },
        {
          href: "https://github.com/6022-labs",
          label: "GitHub",
          position: "right",
        },
      ],
    },
    footer: {
      style: "dark",
      links: [
        {
          title: "Docs",
          items: [
            {
              label: "User Guide",
              to: "/docs/users/intro",
            },
            {
              label: "Developers",
              to: "/docs/developers/intro",
            },
          ],
        },
        {
          title: "Community",
          items: [
            {
              label: "Telegram",
              href: "https://t.me/protocol6022",
            },
            {
              label: "LinkedIn",
              href: "https://www.linkedin.com/company/6022-protocol",
            },
          ],
        },
        {
          title: "More",
          items: [
            {
              label: "GitHub",
              href: "https://github.com/6022-labs/documentation",
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} 6022, Inc. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
