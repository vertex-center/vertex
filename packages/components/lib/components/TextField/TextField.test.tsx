import "@testing-library/jest-dom";
import { fireEvent, render, screen } from "@testing-library/react";
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
            divProps={{ className: "field" }}
        />,
    );
    const field = screen.getByTestId("field");
    expect(field.parentElement).toHaveClass("field", "input");
});

test("it can have a custom input class", () => {
    render(
        <TextField
            id="id"
            data-testid="field"
            placeholder="Placeholder"
            className="input-test"
        />,
    );
    const input = screen.getByTestId("field");
    expect(input).toHaveClass("input-field", "input-test");
});

test("it can react to changes", () => {
    const onChange = jest.fn();
    render(
        <TextField
            id="id"
            data-testid="field"
            placeholder="Placeholder"
            onChange={onChange}
        />,
    );
    const input = screen.getByTestId("field");
    fireEvent.input(input, { target: { value: "New value" } });
    expect(onChange).toHaveBeenCalledTimes(1);
    expect(input).toHaveValue("New value");
});
