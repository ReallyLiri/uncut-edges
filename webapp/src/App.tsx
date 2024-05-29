import styled from "@emotion/styled";
import { Theme } from "./theme.ts";
import { useState } from "react";
import {
  ParserPlaceholder,
  ParserTitle,
  ParserType,
  ParserTypes,
} from "./parser/parserType.ts";
import { parseAsync } from "./parser/api.ts";

const Page = styled.div`
  display: flex;
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
  height: 100vh;
  width: 25vw;
`;

const Row = styled.div`
  display: flex;
  flex-direction: row;
  gap: 1rem;
  padding: 0.25rem 0;
`;

const Button = styled.div<{ disabled: boolean }>`
  background-color: ${() => Theme.White};
  color: ${() => Theme.Black};
  padding: 1rem;
  border-radius: 0.25rem;
  cursor: ${({ disabled }) => (disabled ? "not-allowed" : "pointer")};
  &:hover {
    box-shadow: inset 0 0 0.5rem ${() => Theme.Black};
  }
`;

const StyledInput = styled.input`
  background-color: ${() => Theme.White};
  color: ${() => Theme.Black};
  padding: 0.5rem;
  border-radius: 0.25rem;
  width: 100%;
  font-size: 1rem;
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
  const [loading, setLoading] = useState<boolean>(false);

  const dispatchParsing = async () => {
    try {
      setLoading(true);
      await parseAsync(parserType, value!);
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Page>
      <Container>
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
        {loading ? (
          <div>Working...</div>
        ) : (
          <Button onClick={() => value && dispatchParsing()} disabled={!value}>
            Parse
          </Button>
        )}
      </Container>
    </Page>
  );
}

export default App;
