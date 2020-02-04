#version 100

// out vec4 FragColor;
precision mediump float;
uniform vec4 FragColor;

void main()
{
    gl_FragColor = vec4(1.0, 0.5, 0.2, 1.0);
    // FragColor = vec4(1.0, 0.5, 0.2, 1.0);
    // FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
} 