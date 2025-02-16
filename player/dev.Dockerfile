FROM node:22

WORKDIR /app

COPY . .
RUN npm install

ENV NODE_ENV=development
ENV PORT=3000

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]