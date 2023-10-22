import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import { Input } from "./Input.tsx";
import { createRef } from "react";

test("it renders", () => {
    render(<Input placeholder="Placeholder" />);
    const input = screen.getByPlaceholderText("Placeholder");
    expect(input).toBeInTheDocument();
});

test("it can have a custom class", () => {
    render(<Input placeholder="Placeholder" className="custom-class" />);
    const input = screen.getByPlaceholderText("Placeholder");
    expect(input).toHaveClass("custom-class", "input");
});

test("it can be disabled", () => {
    render(<Input placeholder="Placeholder" disabled />);
    const input = screen.getByPlaceholderText("Placeholder");
    expect(input).toBeDisabled();
});

test("it can be referenced", () => {
    const ref = createRef<HTMLInputElement>();
    render(<Input ref={ref} placeholder="Placeholder" />);
    const input = screen.getByPlaceholderText("Placeholder");
    expect(ref.current).toBe(input);
});
