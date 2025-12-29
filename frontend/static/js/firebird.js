import * as THREE from 'https://cdn.skypack.dev/three@0.160.0';

const canvas = document.getElementById('firebird-canvas');

const renderer = new THREE.WebGLRenderer({
  canvas,
  alpha: true,
  antialias: true
});

const scene = new THREE.Scene();
const camera = new THREE.OrthographicCamera(-1, 1, 1, -1, 0, 1);

const geometry = new THREE.PlaneGeometry(2, 2);

const material = new THREE.ShaderMaterial({
  transparent: true,
  uniforms: {
    uTime: { value: 0 },
    uResolution: { value: new THREE.Vector2() }
  },
  fragmentShader: await fetch('/static/shaders/firebird.frag').then(r => r.text())
});

const mesh = new THREE.Mesh(geometry, material);
scene.add(mesh);

function resize() {
  const rect = canvas.parentElement.getBoundingClientRect();
  renderer.setSize(rect.width, rect.height);
  material.uniforms.uResolution.value.set(rect.width, rect.height);
}

window.addEventListener('resize', resize);
resize();

const start = performance.now();

function animate() {
  material.uniforms.uTime.value = (performance.now() - start) * 0.001;
  renderer.render(scene, camera);
  requestAnimationFrame(animate);
}

animate();

