<script lang="ts">
import { T, useTask, useThrelte } from '@threlte/core';
import {
  DataTexture,
  GLSL3,
  LinearFilter,
  RedFormat,
  RepeatWrapping,
  ShaderMaterial,
  UnsignedByteType,
  Vector2
} from 'three';
import vertexShader from '$lib/shaders/vertex.glsl?raw';

type Props = {
  vertexShader?: string;
  fragmentShader?: string;
};

let {
  fragmentShader
}: Props = $props();

function rand(seed: number) {
  const x = Math.sin(seed) * 43758.5453;
  return x - Math.floor(x);
}

function makeVoronoiTexture(size: number, seed: number) {
  const pointCount = 12;
  const points = Array.from({ length: pointCount }, (_, i) => ({
    x: rand(seed * 13 + i) * size,
    y: rand(seed * 17 + i * 1.3) * size
  }));
  const data = new Uint8Array(size * size);
  for (let y = 0; y < size; y++) {
    for (let x = 0; x < size; x++) {
      let minDist = Infinity;
      for (const p of points) {
        const d = Math.hypot(x - p.x, y - p.y);
        if (d < minDist) minDist = d;
      }
      const v = Math.max(0, Math.min(255, Math.floor((minDist / (size / 4)) * 255)));
      data[y * size + x] = v;
    }
  }
  const texture = new DataTexture(data, size, size, RedFormat, UnsignedByteType);
  texture.needsUpdate = true;
  texture.wrapS = RepeatWrapping;
  texture.wrapT = RepeatWrapping;
  texture.minFilter = LinearFilter;
  texture.magFilter = LinearFilter;
  texture.generateMipmaps = false;
  return texture;
}

const noiseTextureA = makeVoronoiTexture(128, 1);
const noiseTextureB = makeVoronoiTexture(128, 2);
const emptyAudio = new DataTexture(new Uint8Array([0]), 1, 1, RedFormat, UnsignedByteType);
emptyAudio.needsUpdate = true;

const uniforms = {
  u_time: { value: 500 },
  u_resolution: { value: new Vector2(1, 1) },
  u_mouse: { value: new Vector2(0, 0) },
  iChannel0: { value: noiseTextureA },
  iChannel1: { value: noiseTextureB },
  u_noiseA: { value: noiseTextureA },
  u_noiseB: { value: noiseTextureB }
};

const threlte = useThrelte();

useTask((delta) => {
  const { width, height } = threlte.size.current;
  const dpr = threlte.dpr.current;
  uniforms.u_resolution.value.set(width * dpr, height * dpr);
  uniforms.u_time.value += delta * 1;
});
</script>

<T.Mesh position={[0, 0, 0]}>
  <T.PlaneGeometry args={[30, 10]} />
  <T.ShaderMaterial
    glslVersion={GLSL3}
    fragmentShader={fragmentShader}
    {uniforms}
  />
</T.Mesh>

  