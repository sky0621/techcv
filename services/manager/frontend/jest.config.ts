import type { Config } from 'jest';

const config: Config = {
  preset: 'ts-jest',
  testEnvironment: 'jsdom',
  moduleDirectories: ['node_modules', 'src'],
  moduleNameMapper: {
    '^@/lib/env$': '<rootDir>/src/lib/env.jest.ts',
    '^@/(.*)$': '<rootDir>/src/$1',
    '\\.(css|less|scss|sass)$': '<rootDir>/testing/__mocks__/styleMock.ts'
  },
  setupFilesAfterEnv: ['<rootDir>/testing/setup-tests.ts'],
  testPathIgnorePatterns: ['/node_modules/', '/dist/'],
  coveragePathIgnorePatterns: ['/node_modules/', '/testing/'],
  globals: {
    'ts-jest': {
      tsconfig: '<rootDir>/tsconfig.json'
    }
  }
};

export default config;
