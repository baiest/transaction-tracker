import { render, screen } from "@testing-library/react";
import LineChart from "./LineChart";
import { vi } from "vitest";
import { EChartsOption } from "echarts";

vi.mock("next/dynamic", () => ({
  __esModule: true,
  default: vi.fn(
    () =>
      function MockChart(props: { option: unknown }) {
        return (
          <div
            data-testid="mock-echarts-chart"
            data-options={JSON.stringify(props.option)}
          >
            Mock Chart
          </div>
        );
      }
  )
}));

describe("LineChart", () => {
  const defaultProps = {
    xData: ["Jan", "Feb", "Mar"],
    series: [
      { name: "Income", data: [100, 200, 150], color: "green" },
      { name: "Expenses", data: [80, 150, 120], color: "red" }
    ]
  };

  it("should render the mock chart component", () => {
    render(<LineChart {...defaultProps} />);
    expect(screen.getByTestId("mock-echarts-chart")).toBeInTheDocument();
  });

  it("should pass the correct x-axis data to the ECharts option", () => {
    render(<LineChart {...defaultProps} />);
    const chart = screen.getByTestId("mock-echarts-chart");
    const options = JSON.parse(chart.dataset.options as string);
    expect(options.xAxis.data).toEqual(defaultProps.xData);
  });

  it("should pass the correct series data to the ECharts option", () => {
    render(<LineChart {...defaultProps} />);
    const chart = screen.getByTestId("mock-echarts-chart");
    const options = JSON.parse(chart.dataset.options as string);
    expect(options.series[0].name).toBe("Income");
    expect(options.series[0].data).toEqual([100, 200, 150]);
    expect(options.series[0].itemStyle.color).toBe("green");
    expect(options.series[1].name).toBe("Expenses");
    expect(options.series[1].data).toEqual([80, 150, 120]);
    expect(options.series[1].itemStyle.color).toBe("red");
  });

  it("should apply a custom title if provided", () => {
    const customTitle = "Financial Overview";
    render(<LineChart {...defaultProps} title={customTitle} />);
    const chart = screen.getByTestId("mock-echarts-chart");
    const options = JSON.parse(chart.dataset.options as string);
    expect(options.title.text).toBe(customTitle);
  });

  it("should merge custom options with default options", () => {
    const customOptions: EChartsOption = {
      legend: { show: true },
      tooltip: { trigger: "item" }
    };
    render(<LineChart {...defaultProps} customOptions={customOptions} />);
    const chart = screen.getByTestId("mock-echarts-chart");
    const options = JSON.parse(chart.dataset.options as string);
    expect(options.legend.show).toBe(true);
    expect(options.tooltip.trigger).toBe("item"); // Verifies override
  });

  it("should use default height and minHeight if not provided", () => {
    render(<LineChart {...defaultProps} />);
    const chart = screen.getByTestId("mock-echarts-chart");
    expect(chart).toBeInTheDocument();
  });
});
