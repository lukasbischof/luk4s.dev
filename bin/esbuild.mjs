#!/usr/bin/env node

import { sassPlugin } from "esbuild-sass-plugin";
import esbuild from "esbuild";
import autoprefixer from "autoprefixer";
import postcss from "postcss";
import postcssPresetEnv from "postcss-preset-env";
import minify from "postcss-minify";

await esbuild.build({
    entryPoints: ["assets/entrypoints/application.scss"],
    bundle: true,
    outdir: "public/build",
    assetNames: "images/[name]",
    publicPath: "",
    loader: {
        ".jpeg": "file",
        ".svg": "file",
    },
    plugins: [sassPlugin({
        transform: async (source) => {
            const processor = postcss([
                autoprefixer,
                postcssPresetEnv({ stage: 0 }),
                minify,
            ]);
            const { css } = await processor.process(source, { from: undefined });
            console.log(css);
            return css;
        },
        // precompile(source, pathname) {
        //     console.log(pathname)
        //     const basedir = path.dirname(pathname)
        //     console.log(basedir)
        //     return source.replace(/(url\(['"]?)(\.\.?\/)([^'")]+['"]?\))/g, `$1${basedir}/$2$3`)
        // }
    })],
});
