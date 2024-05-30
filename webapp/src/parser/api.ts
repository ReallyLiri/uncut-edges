import { ParserType } from "./parserType.ts";
import { WritableStream } from "web-streams-polyfill";
import streamSaver from "streamsaver";

if (!window.WritableStream) {
  // @ts-ignore
  streamSaver.WritableStream = WritableStream;
  // @ts-ignore
  window.WritableStream = WritableStream;
}

const API_URL =
  import.meta.env.VITE_API_URL || "https://uncut-edges.onrender.com";

const ParserRoute: Record<ParserType, string> = {
  manifest: "/parse/",
  penn: "/parse/penn/",
  shakespeare: "/parse/shakespeare/",
};

export const parseAsync = (
  type: ParserType,
  input: string,
  pageRange: string | undefined,
): Promise<void> => {
  let url = `${API_URL}${ParserRoute[type]}${encodeURIComponent(input)}`;
  if (pageRange) {
    url += `?pages=${encodeURIComponent(pageRange)}`;
  }
  console.log("calling parser", url);
  return new Promise<void>((resolve, reject) => {
    fetch(url)
      .then((response) => {
        const contentDisposition = response.headers.get("Content-Disposition")!;
        const fileName = contentDisposition.substring(
          contentDisposition.lastIndexOf("=") + 1,
        );

        const fileStream = streamSaver.createWriteStream(fileName);
        const readableStream = response.body!;
        if (readableStream.pipeTo) {
          return readableStream.pipeTo(fileStream).then(resolve).catch(reject);
        }
        const writer = fileStream.getWriter();

        const reader = readableStream.getReader();

        const pump: () => void = () =>
          reader
            .read()
            .then((res) =>
              res.done
                ? writer.close().then(resolve).catch(reject)
                : writer.write(res.value).then(pump).catch(reject),
            )
            .catch(reject);

        pump();
      })
      .catch(reject);
  });
};
