import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import { TextField } from "./TextField.tsx";

test("it renders", () => {
    render(<TextField data-testid="field" placeholder="Placeholder" />);
    const field = screen.getByTestId("field");
    expect(field).toBeInTheDocument();
});

test("it can have a custom class", () => {
    render(
        <TextField
            data-testid="field"
            placeholder="Placeholder"
            className="field"
        />,
    );
    const field = screen.getByTestId("field");
    expect(field).toHaveClass("field", "text-field");
});

test("it can have a custom input class", () => {
    render(
        <TextField
            data-testid="field"
            placeholder="Placeholder"
            inputProps={{ className: "input" }}
        />,
    );
    const input = screen.getByTestId("field").children[1];
    expect(input).toHaveClass("input");
});
