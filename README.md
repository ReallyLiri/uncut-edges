# Uncut Edges

iiif to PDF tool , still WIP

## WebApp

Availble at [https://uncut-edges.netlify.app](https://uncut-edges.netlify.app)

[![Netlify Status](https://api.netlify.com/api/v1/badges/21eab1bb-8b18-4059-8a40-b93ac78b4184/deploy-status)](https://app.netlify.com/sites/uncut-edges/deploys)

## API

Public API available at `https://uncut-edges.onrender.com`

### Endpoints

* `GET /parse/<manifest URL>` - Parse an iiif manifest by URL, i.e `/parse/https%3A%2F%2Fcolenda.library.upenn.edu%2Fitems%2Fark%3A%2F81431%2Fp3hk28%2Fmanifest`
* `GET /parse/penn/<catalog ID>` - Parse from a catalog in Colenda Digital Repository of Penn Libraries, i.e `/parse/penn/81431-p3hk28`
* `GET /parse/shakespeare/<catalog ID>` - Parse from a catalog in Shakespeare Digital Repository of Penn Libraries, i.e `/parse/shakespeare/bib244741-309974-lb41`
