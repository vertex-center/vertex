import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import Button from "./index.tsx";

test("it renders", () => {
    render(<Button>Button</Button>);
    const button = screen.getByRole("button");
    expect(button).toBeInTheDocument();
});

test("it can be clicked", () => {
    const onClick = jest.fn();
    render(<Button onClick={onClick}>Button</Button>);
    const button = screen.getByRole("button");
    button.click();
    expect(onClick).toHaveBeenCalled();
});

test("it can be disabled", () => {
    render(<Button disabled>Button</Button>);
    const button = screen.getByRole("button");
    expect(button).toBeDisabled();
    expect(button).toHaveClass("button-disabled");
});

test("it can have a custom class", () => {
    render(<Button className="custom-class">Button</Button>);
    const button = screen.getByRole("button");
    expect(button).toHaveClass("custom-class", "button");
});

test("it is outlined by default", () => {
    render(<Button>Button</Button>);
    const button = screen.getByRole("button");
    expect(button).toHaveClass("button-outlined");
});

describe("it can be of type", () => {
    test("colored", () => {
        render(<Button type="colored">Button</Button>);
        const button = screen.getByRole("button");
        expect(button).toHaveClass("button-colored");
    });

    test("outlined", () => {
        render(<Button type="outlined">Button</Button>);
        const button = screen.getByRole("button");
        expect(button).toHaveClass("button-outlined");
    });
});
