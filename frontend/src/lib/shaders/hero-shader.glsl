precision mediump float;
layout(location = 0) out vec4 pc_fragColor;

// Shadertoy uniforms (mapped from WebGL uniforms)
uniform float u_time;
uniform vec2 u_resolution;
uniform vec2 u_mouse;
uniform float u_audioAmplitude;
uniform sampler2D u_audioMap;
uniform vec2 u_audioMapSize;
uniform float u_audioMapFps;
uniform float u_audioTime;

// Shadertoy texture channels
uniform sampler2D iChannel0;
uniform sampler2D iChannel1;

// Shadertoy variables (aliases for compatibility)
#define iTime u_time
#define iResolution u_resolution
#define iMouse u_mouse

#define PI 3.14159265358979
#define SPEED 0.05
#define GRID_COUNT 64.
#define STEP 0.04
#define COUNT 96.
#define HEIGHT 0.1
#define SUN_SPEED 0.1
#define MIN_HEIGHT 0.01
#define MAX_HEIGHT 0.5

vec3 GetBackgroundColor(in vec3 ray)
{
    vec3 c0 = vec3(115. / 255., 25. / 255., 171. / 255.);
    vec3 c1 = vec3(235. / 255., 51. / 255., 201. / 255.);
    vec3 c2 = vec3(40. / 255., 7. / 255., 80. / 255.);
    float ma = 1.0;
    float a = fract(iTime * 0.25) * ma;
    a = a > ma/2.0 ? ma - a : a;
    vec3 c = c0 * a + c2 * (1.0 - a);
    a = inversesqrt(max(1e-3, 1.0 - abs(ray.y)));
    a = clamp(a * 0.7, 0.0, 1.0);
    float b = abs(ray.x * 0.5);
    float m = texture(iChannel1, vec2(b * 0.02, fract(iTime * 0.00002))).x;
    a = 0.05 * m;
        
    c = c1 * a + c * (1.0 - a);
    return vec3(c);
}

bool GetHitWithPlane(in vec3 o, in vec3 ray, out vec3 hitPoint)
{
    bool hit = false;
    hitPoint = vec3(0.0, 0.0, 0.0);
    float realHeight = 0.;
    float preRealHeight = 0.;
    float _step = STEP;// min(STEP / abs(ray.y), 2. * STEP);
    for (float i = 1.; i <= COUNT; i++)
    {
        vec3 current = o + ray * _step * i;
        realHeight = texture(iChannel0, current.xz * 0.5).x * HEIGHT * .4;
        if (realHeight > current.y + 0.2) {
            hitPoint = current;
            return true;
        }
        
        float s = abs(current.x - o.x) / 0.5;
        s = clamp(s * s, 0., 1.);
        realHeight *= s;
        if (current.y < realHeight)
        {
            vec3 last = current - ray * _step;
            float preDiff = last.y - preRealHeight;
            float curDiff = realHeight - current.y;
            float t = preDiff / (preDiff + curDiff);
            current = last + ray * _step * t;      
            hitPoint = current;
            return true;
        }
        preRealHeight = realHeight;
    }
    return hit;
}

vec3 GetHitPointColor(in vec3 cameraPos, in vec3 hitPoint, in vec3 background)
{
    float fx = abs(fract(hitPoint.x * GRID_COUNT + 0.5) - 0.5);
    float fz = abs(fract(hitPoint.z * GRID_COUNT + 0.5) - 0.5);

    vec3 b = vec3(0.2, 0.6, 0.9);
    vec3 r = background;
    //r.x = pow(r.x, 3.3);
    //r.y = pow(r.y, 0.7);
    //r.z = pow(r.z, 0.5);
    float d = length(hitPoint.xz - cameraPos.xz) * 1.0;
    d = clamp(d, 0., 1.);
    vec3 col = r * d + b * (1. - d);
    float li = texture(iChannel0, hitPoint.xz / 64. + vec2(0, -iTime * 0.005)).x * 0.5 + 0.5;
    float lo = abs(hitPoint.z - cameraPos.z) + 0.3;
    col /= min(li, lo);
    
    
    d = min(fx, fz);
    float bloom = smoothstep(d, 0., 0.1);
    d = smoothstep(d, 0., 0.01);
    d += bloom * 0.2;
    
    return r * d + col * (1. - d);
}

vec4 GetSun(in vec3 cameraPos, in vec3 ray)
{
    if (dot(ray, normalize(vec3(0.0, 0.035, 1.0))) > 0.97)
    {
        vec4 r = vec4(1., 0., 0., 1.);
        vec4 y = vec4(2., 2., 0., 1.);
        float a = abs(ray.y);
        vec4 col = y * a + r * (1. - a);
        
        float m = 0.5;
        float g = 0.2;
        float ft = fract(iTime * SUN_SPEED) * m;
        vec2 l0 = vec2(-ft + g, ft * 0.2);
        ft = fract(iTime * SUN_SPEED + 0.25) * m;
        vec2 l1 = vec2(-ft + g, ft * 0.2);
        ft = fract(iTime * SUN_SPEED + 0.5) * m;
        vec2 l2 = vec2(-ft + g, ft * 0.2);
        ft = fract(iTime * SUN_SPEED + 0.75) * m;
        vec2 l3 = vec2(-ft + g, ft * 0.2);
        return col * (1. - 
                      step(l0.x - l0.y, ray.y) * step(ray.y, l0.x + l0.y) -
                      step(l1.x - l1.y, ray.y) * step(ray.y, l1.x + l1.y) -
                      step(l2.x - l2.y, ray.y) * step(ray.y, l2.x + l2.y) -
                      step(l3.x - l3.y, ray.y) * step(ray.y, l3.x + l3.y)
                      );
        
        return col * 
                (step(ray.y, -0.03) + 
                 step(0., ray.y) * step(ray.y, 0.04) +
                 step(0.06, ray.y) * step(ray.y, 0.11) +
                 step(0.13, ray.y) * step(ray.y, 0.18) +
                 step(0.195, ray.y)
                 );
    }
    return vec4(0.0, 0.0, 0.0, 0.0);
}

vec4 GetScreenColor(in vec3 ray, in vec3 cameraPos)
{
    vec3 background = GetBackgroundColor(ray);
    bool hit;
    vec3 hitPoint;
    hit = GetHitWithPlane(cameraPos, ray, hitPoint);
    vec3 plane = GetHitPointColor(cameraPos, hitPoint, background);
    if (hit)
        return vec4(plane, 1);
    vec4 sun = GetSun(cameraPos, ray);
    
    return vec4(background * (1.0 - sun.w) + sun.xyz * sun.w, 1.);
}

void mainImage( out vec4 fragColor, in vec2 fragCoord )
{
    vec2 uv = fragCoord/iResolution.xy;
    vec2 cuv = uv * 2.0 - 1.0;
    cuv.x *= iResolution.x / iResolution.y;
    vec3 ray = normalize(vec3(cuv.x, cuv.y, 1.2));
    
    float ix = PI * 0.5 + (iMouse.x / iResolution.x) * PI;
    float iy = iMouse.y / iResolution.y * PI * 0.25;
    
    vec3 f = vec3(
        cos(ix) * cos(iy),
        sin(iy),
        sin(ix) * cos(iy)
    );
    f = normalize(f);
    vec3 u = vec3(0.0, 1.0, 0.0);
    vec3 r = normalize(cross(f, u));
    u = normalize(cross(r, f));
    
    
    fragColor = GetScreenColor(normalize(f * ray.z + u * ray.y + r * ray.x), vec3(0.0, 0.1, iTime * SPEED));
}

// Standard GLSL main function that calls mainImage
void main() {
	vec2 fragCoord = gl_FragCoord.xy;
	vec4 color;
	mainImage(color, fragCoord);
	pc_fragColor = color;
}
