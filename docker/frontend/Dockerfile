FROM node:alpine AS builder

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

RUN npm run build 

# Run TypeScript application in dev
# CMD npm run dev


# Stage 2: Serve with Nginx
FROM nginx:latest
COPY --from=builder /app/dist /usr/share/nginx/html
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
