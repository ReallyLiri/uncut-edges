import "./Game.css";

import Resources from "./Resources.js";
import { useEffect, useRef } from "react";
import { DinoScript } from "./DinoScript.ts";

const appendScript = (toElement: HTMLDivElement, script: string) => {
  const dinoScriptContainer = document.createElement("script");
  dinoScriptContainer.appendChild(document.createTextNode(script));
  toElement.appendChild(dinoScriptContainer);
};

let loaded = false;

export const Game = () => {
  const startDiv = useRef<HTMLDivElement>(null);
  const endDiv = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (loaded || !startDiv.current || !endDiv.current) {
      return;
    }
    loaded = true;
    appendScript(startDiv.current, DinoScript);
    appendScript(endDiv.current, "new Runner('.interstitial-wrapper');");
    return () => {
      loaded = false;
      document.querySelector(".runner-container")?.remove();
    };
  }, [startDiv.current, endDiv.current]);

  return (
    <div ref={startDiv}>
      <div className="interstitial-wrapper">
        <Resources />
        <div ref={endDiv}></div>
      </div>
    </div>
  );
};
