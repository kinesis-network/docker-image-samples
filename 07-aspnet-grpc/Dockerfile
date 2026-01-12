ARG DOTNET_VERSION=9.0
ARG PUBDIR=/app/publish
ARG SRCDIR=/app/src

FROM mcr.microsoft.com/dotnet/sdk:${DOTNET_VERSION} AS build
WORKDIR /app

ARG DOTNET_CONFIG=Release
ARG PUBDIR
ARG SRCDIR

RUN mkdir -p ${SRCDIR}
COPY . ${SRCDIR}/

RUN cd ${SRCDIR} \
  && dotnet publish \
    -c ${DOTNET_CONFIG} \
    -p:UseAppHost=false \
    -o ${PUBDIR}

FROM mcr.microsoft.com/dotnet/aspnet:${DOTNET_VERSION}
WORKDIR /app

ARG PUBDIR
ARG SRCDIR

RUN apt update -y && apt install -y --no-install-recommends \
  curl \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

COPY --from=build ${PUBDIR} ${PUBDIR}
COPY --chmod=775 entrypoint.sh .

ENV ASPNETCORE_URLS=http://0.0.0.0:5000
ENTRYPOINT ["/app/entrypoint.sh"]
