# Use the official Node.js v16.x LTS image as the base image
FROM node:16-alpine3.14

# Set the working directory to /app
WORKDIR /app

# Copy package.json and package-lock.json files to the container
COPY ./frontend .

# Install the project dependencies
RUN npm install

# Build the project for production
RUN npm run build

# Expose port 3002
EXPOSE 5173

# Start the application
CMD ["npm", "run", "dev"]