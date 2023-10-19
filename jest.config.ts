import { JestConfigWithTsJest } from "ts-jest";

const config: JestConfigWithTsJest = {
    transform: {
        "^.+\\.svg$": "jest-transform-stub",
    },
    preset: "ts-jest",
    testEnvironment: "jsdom",
    moduleNameMapper: {
        "^@/(.*)$": "<rootDir>/lib/$1",
        "^.+\\.(sass)$": "identity-obj-proxy",
    },
};

export default config;
