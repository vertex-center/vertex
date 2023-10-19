import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import Button from "./index.tsx";
import { ButtonType } from "./Button.tsx";
import { MaterialIcon } from "../../index.ts";

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

describe("it can be of variants", () => {
    const cases: ButtonType[] = ["colored", "outlined", "danger"];
    test.each(cases)("%p", (type) => {
        render(<Button variant={type}>Button</Button>);
        const button = screen.getByRole("button");
        expect(button).toHaveClass(`button-${type}`);
    });
});

describe("it can have an icon", () => {
    const icon = <MaterialIcon icon="deployed_code_update" />;
    render(
        <Button leftIcon={icon} rightIcon={icon}>
            Button
        </Button>,
    );
    const buttons = screen.getByRole("button");
    expect(buttons).toHaveLength(2);
    expect(buttons).toContainElement(screen.getByText("deployed_code_update"));
});
