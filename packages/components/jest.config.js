const config = {
    transform: {
        "^.+\\.svg$": "jest-transform-stub",
    },
    preset: "ts-jest",
    testEnvironment: "jsdom",
    moduleNameMapper: {
        "^@/(.*)$": "<rootDir>/lib/$1",
        "^.+\\.(sass|css)$": "identity-obj-proxy",
    },
};

export default config;
