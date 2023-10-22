import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import Input from "./Input.tsx";

test("it renders", () => {
    render(<Input placeholder="Placeholder" />);
    const button = screen.getByPlaceholderText("Placeholder");
    expect(button).toBeInTheDocument();
});

test("it can have a custom class", () => {
    render(<Input placeholder="Placeholder" className="custom-class" />);
    const button = screen.getByPlaceholderText("Placeholder");
    expect(button).toHaveClass("custom-class", "input");
});

test("it can be disabled", () => {
    render(<Input placeholder="Placeholder" disabled />);
    const button = screen.getByPlaceholderText("Placeholder");
    expect(button).toBeDisabled();
});
