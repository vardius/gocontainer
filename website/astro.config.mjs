// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
	site: 'https://vardius.github.io',
	base: '/gocontainer',
	trailingSlash: 'always',
	// Pages import the runnable programs in ../examples, so the dev server may read outside this folder.
	vite: { server: { fs: { allow: ['..'] } } },
	integrations: [
		starlight({
			title: 'gocontainer',
			description:
				'A dependency injection container for Go with typed keys: register services and get them back as their own types, with no type assertions or reflection.',
			logo: { src: './src/assets/mark.svg', alt: '' },
			favicon: '/favicon.svg',
			components: { Hero: './src/components/Hero.astro' },
			// Long lines wrap instead of scrolling under the copy button on phones.
			expressiveCode: { defaultProps: { wrap: true } },
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/vardius/gocontainer' }],
			customCss: ['./src/styles/custom.css'],
			sidebar: [
				{
					label: 'Start here',
					items: ['getting-started'],
				},
				{
					label: 'Guides',
					items: ['keys', 'container', 'testing'],
				},
				{
					label: 'Reference',
					items: [
						'benchmarks',
						'migrating-from-v1',
						{ label: 'API reference', link: 'https://pkg.go.dev/github.com/vardius/gocontainer/v2' },
					],
				},
			],
		}),
	],
});
