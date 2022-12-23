#!/usr/bin/env node

import autoprefixer from "autoprefixer";
import esbuild from "esbuild";
import minify from "postcss-minify";
import postcss from "postcss";
import postcssPresetEnv from "postcss-preset-env";
import { sassPlugin } from "esbuild-sass-plugin";

await esbuild.build({
    assetNames: "images/[name]",
    bundle: true,
    entryPoints: ["assets/entrypoints/application.scss"],
    loader: { ".jpeg": "file", ".svg": "file" },
    minify: true,
    outdir: "public/build",
    publicPath: "",
    plugins: [sassPlugin({
        transform: async (source) => {
            const processor = postcss([
                autoprefixer,
                postcssPresetEnv({ stage: 0 }),
                minify,
            ]);
            return (await processor.process(source, { from: undefined })).css;
        },
    })],
});
