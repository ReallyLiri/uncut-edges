# Uncut Edges

iiif to PDF tool , still WIP

## API

Public API available at `https://uncut-edges.onrender.com`

### Endpoints

* `GET /parse/<manifest URL>` - Parse an iiif manifest by URL, i.e `/parse/https%3A%2F%2Fcolenda.library.upenn.edu%2Fitems%2Fark%3A%2F81431%2Fp3hk28%2Fmanifest`
* `GET /parse/penn/<catalog ID>` - Parse from a catalog in Colenda Digital Repository of Penn Libraries, i.e `/parse/penn/81431-p3hk28`
* `GET /parse/shakespeare/<catalog ID>` - Parse from a catalog in Shakespeare Digital Repository of Penn Libraries, i.e `/parse/shakespeare/bib244741-309974-lb41`
