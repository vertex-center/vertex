import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import { Input } from "./Input.tsx";
import { createRef } from "react";

test("it renders", () => {
    render(<Input id="id" placeholder="Placeholder" />);
    const input = screen.getByPlaceholderText("Placeholder");
    expect(input).toBeInTheDocument();
});

test("it can have a custom class", () => {
    render(
        <Input
            id="id"
            data-testid="input"
            placeholder="Placeholder"
            className="custom-class"
        />,
    );
    const input = screen.getByTestId("input");
    expect(input).toHaveClass("custom-class", "input");
});

test("it can be disabled", () => {
    render(<Input id="id" placeholder="Placeholder" disabled />);
    const input = screen.getByPlaceholderText("Placeholder");
    expect(input).toBeDisabled();
});

test("it can be referenced", () => {
    const ref = createRef<HTMLInputElement>();
    render(<Input id="id" inputProps={{ ref }} placeholder="Placeholder" />);
    const input = screen.getByPlaceholderText("Placeholder");
    expect(ref.current).toBe(input);
});
