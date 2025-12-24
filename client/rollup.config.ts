import typescript from '@rollup/plugin-typescript';
import type { RollupOptions } from 'rollup'

export default {
  input: 'src/index.ts',
  output: [
    {
      file: 'dist/index.js',
      format: 'cjs',
      sourcemap: true,
    },
    {
      file: 'dist/index.esm.js',
      format: 'esm',
      sourcemap: true,
    },
    {
      file: 'dist/index.umd.js',
      format: 'umd',
      name: 'Zanalytics',
      sourcemap: true,
    }
  ],
  plugins: [
    typescript({
      tsconfig: './tsconfig.json'
    })
  ]
} as RollupOptions
