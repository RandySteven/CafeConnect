FROM node:24.1.0-alpine

WORKDIR /app

COPY . .

COPY package.json package-lock.json* ./

RUN npm install

ENV PORT 3000

EXPOSE 3000

CMD ["npm", "start"]
