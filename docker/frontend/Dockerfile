FROM node:alpine AS build

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json from frontend/ to the working directory
COPY frontend/package*.json ./

# Copy tsconfig and vite.config to resolve paths for shadcn-ui
COPY frontend/components.json ./
COPY frontend/tsconfig.json ./
COPY frontend/vite.config.ts ./

# import "path" without error for shadcn-ui
RUN npm i -D @types/node

# Install dependencies
RUN npm install

# Copy the rest of the frontend application source code to the working directory
COPY frontend/ .

# Run TypeScript application in dev
CMD npm run dev
