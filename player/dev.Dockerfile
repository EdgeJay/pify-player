FROM node:22

WORKDIR /app

COPY . .

RUN npm install

ENV NODE_ENV=development

EXPOSE 5173
EXPOSE 24678

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]