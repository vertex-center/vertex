import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import { TextField } from "./TextField.tsx";

test("it renders", () => {
    render(<TextField id="id" data-testid="field" placeholder="Placeholder" />);
    const field = screen.getByTestId("field");
    expect(field).toBeInTheDocument();
});

test("it can have a custom class", () => {
    render(
        <TextField
            id="id"
            data-testid="field"
            placeholder="Placeholder"
            className="field"
        />,
    );
    const field = screen.getByTestId("field");
    expect(field).toHaveClass("field", "input");
});

test("it can have a custom input class", () => {
    render(
        <TextField
            id="id"
            data-testid="field"
            placeholder="Placeholder"
            inputProps={{ className: "input" }}
        />,
    );
    const input = screen.getByTestId("field").children[0];
    expect(input).toHaveClass("input-field");
});
