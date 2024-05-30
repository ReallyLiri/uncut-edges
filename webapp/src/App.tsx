import styled from "@emotion/styled";
import { Theme } from "./theme.ts";
import { useEffect, useState } from "react";
import {
  ParserPlaceholder,
  ParserTitle,
  ParserType,
  ParserTypes,
} from "./parser/parserType.ts";
import { parseAsync } from "./parser/api.ts";
import { Game } from "./game/Game.tsx";

const Page = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  background-color: ${() => Theme.Black};
  color: ${() => Theme.White};
`;

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
  justify-content: center;
  align-items: start;
  height: fit-content;
  width: 25vw;
  padding: 2rem;
  border: 1px solid ${() => Theme.White};
  border-radius: 0.5rem;
`;

const MainTitle = styled.div`
  font-size: 4rem;
  font-weight: 600;
  margin-bottom: 2rem;
  font-family: "Jaini", system-ui;
  pointer-events: none;
  user-select: none;
`;

const Title = styled.div`
  font-weight: 600;
`;

const Row = styled.div`
  display: flex;
  flex-direction: row;
  gap: 1rem;
  padding: 0.25rem 0;
  align-items: center;
  width: 100%;
`;

const Button = styled.div<{ disabled?: boolean }>`
  background-color: ${() => Theme.White};
  color: ${() => Theme.Black};
  padding: 1rem;
  border-radius: 0.25rem;
  cursor: ${({ disabled }) => (disabled ? "not-allowed" : "pointer")};

  &:hover {
    box-shadow: inset 0 0 0.5rem ${() => Theme.Black};
  }
`;

const Label = styled.div`
  font-size: 0.75rem;
  opacity: 0.5;
`;

const PlayLabel = styled(Label)`
  padding-top: 0.5rem;
  text-align: start;
  width: calc(25vw + 4rem);
  cursor: pointer;
`;

const StyledInput = styled.input`
  background-color: ${() => Theme.White};
  color: ${() => Theme.Black};
  padding: 0.5rem;
  border-radius: 0.25rem;
  width: -webkit-fill-available;
  font-size: 1rem;
`;

const ErrorMessage = styled.div`
  padding: 1rem;
  font-size: 1.5rem;
  font-weight: 600;
  color: #9e3333;
`;

const TextInput = ({
  onChange,
  placeholder,
}: {
  onChange: (value: string) => void;
  placeholder: string;
}) => (
  <StyledInput
    type="text"
    onChange={(e) => onChange(e.target.value)}
    placeholder={placeholder}
  />
);

const ParserInput = ({
  selected,
  value,
}: {
  selected: ParserType;
  value: ParserType;
}) => (
  <input
    type="radio"
    name="parser"
    readOnly
    value={value}
    checked={selected === value}
  />
);

function App() {
  const [value, setValue] = useState<string>();
  const [parserType, setParserType] = useState<ParserType>("manifest");
  const [pageRange, setPageRange] = useState<string>("");
  const [gameOn, setGameOn] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>();
  const [loadingDots, setLoadingDots] = useState(1);

  const dispatchParsing = async () => {
    try {
      setError(undefined);
      setLoading(true);
      setGameOn(true);
      await parseAsync(parserType, value!, pageRange);
    } catch (e) {
      console.error(e);
      setError((e as any)?.message || "An error occurred");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!loading) {
      return;
    }
    const interval = setInterval(() => {
      setLoadingDots((d) => (d + 1) % 3);
    }, 500);
    return () => clearInterval(interval);
  }, [loading]);

  if (gameOn) {
    return (
      <Page>
        <div>
          {loading
            ? `Working${Array(loadingDots + 2).join(".")} This may take a few minutes. Why not play a game in the meantime?`
            : "Parsing completed!"}
        </div>
        <Game />
        {error && <ErrorMessage>Error: {error}</ErrorMessage>}
        {!loading && (
          <Button onClick={() => setGameOn(false)}>Back to parser</Button>
        )}
      </Page>
    );
  }

  return (
    <Page>
      <MainTitle>Uncut Edges</MainTitle>
      <Container>
        <Title>IIIF Parser 📚</Title>
        <div onChange={(e) => setParserType((e.target as any).value)}>
          {ParserTypes.map((type) => (
            <Row key={type} onClick={() => setParserType(type)}>
              <ParserInput selected={parserType} value={type} />{" "}
              {ParserTitle[type]}
            </Row>
          ))}
        </div>
        <TextInput
          placeholder={ParserPlaceholder[parserType]}
          onChange={setValue}
        />
        <Row>
          <div>
            <Label>Optional</Label> Pages:
          </div>
          <TextInput placeholder="i.e 1-3, 5, 7" onChange={setPageRange} />
        </Row>
        {loading ? (
          <div>Working...</div>
        ) : (
          <Button onClick={() => value && dispatchParsing()} disabled={!value}>
            Parse
          </Button>
        )}
      </Container>
      <PlayLabel onClick={() => setGameOn(true)}>
        or are you just here to play?
      </PlayLabel>
    </Page>
  );
}

export default App;
