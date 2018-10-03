	# Build stage
FROM golang AS build
ADD . /src
RUN apt-get update -yq \
    && apt-get install curl gnupg -yq \
    && curl -sL https://deb.nodesource.com/setup_8.x | bash \
    && apt-get install nodejs -yq
RUN cd /src && CGO_ENABLED=0 go build -o fpg
RUN cd /src/static && npm install && npm run-script build

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build /src/fpg /app/
COPY --from=build /src/static /app/static
ENTRYPOINT ./fpg
EXPOSE 9090
