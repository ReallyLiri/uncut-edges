export type ParserType = "manifest" | "penn" | "shakespeare";

export const ParserPlaceholder: Record<ParserType, string> = {
  manifest: "https://example.com/manifest.json",
  penn: "Catalog ID, i.e 81431-p3hk28",
  shakespeare: "Catalog ID, i.e bib244741-309974-lb41",
};

export const ParserTitle: Record<ParserType, string> = {
  manifest: "Manifest URL",
  penn: "Penn Libraries: Colenda Digital Repository",
  shakespeare: "Folger Shakespeare Library: Digital Collections",
};

export const ParserTypes: ParserType[] = Object.keys(
  ParserTitle,
) as ParserType[];
