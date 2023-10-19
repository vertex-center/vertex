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

test("it can have an icon", () => {
    render(
        <Button
            leftIcon={<MaterialIcon icon="deployed_code_update" />}
            rightIcon={<MaterialIcon icon="arrow_forward" />}
        >
            Button
        </Button>,
    );
    const button = screen.getByRole("button");
    expect(button.children).toHaveLength(3);
    expect(button).toContainElement(screen.getByText("deployed_code_update"));
    expect(button).toContainElement(screen.getByText("arrow_forward"));
});
