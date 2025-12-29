precision highp float;

uniform float uTime;
uniform vec2 uResolution;

float noise(vec2 p) {
    return fract(sin(dot(p, vec2(127.1, 311.7))) * 43758.5453);
}

float fbm(vec2 p) {
    float v = 0.0;
    float a = 0.5;
    for (int i = 0; i < 5; i++) {
        v += a * noise(p);
        p *= 2.0;
        a *= 0.5;
    }
    return v;
}

// Very rough phoenix silhouette
float birdSDF(vec2 uv) {
    uv.x *= 1.4;

    float body = length(uv) - 0.25;
    float wingL = length(uv + vec2(0.35, 0.1)) - 0.35;
    float wingR = length(uv - vec2(0.35, 0.1)) - 0.35;

    return min(body, min(wingL, wingR));
}

void main() {
    vec2 uv = (gl_FragCoord.xy / uResolution) * 2.0 - 1.0;
    uv.y *= uResolution.y / uResolution.x;

    float sdf = birdSDF(uv);
    float mask = smoothstep(0.02, -0.02, sdf);

    vec2 flameUV = uv;
    flameUV.y += uTime * 0.4;

    float fire = fbm(flameUV * 3.0);
    fire = pow(fire, 2.0);

    vec3 color = mix(
        vec3(1.0, 0.3, 0.0),
        vec3(1.0, 1.0, 0.2),
        fire
    );

    gl_FragColor = vec4(color * mask, mask);
}

