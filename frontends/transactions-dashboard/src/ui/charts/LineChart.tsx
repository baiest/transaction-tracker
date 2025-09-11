"use client";

import dynamic from "next/dynamic";
import { FC } from "react";
import type { EChartsOption } from "echarts";
import type { EChartsReactProps } from "echarts-for-react";

const ReactECharts = dynamic<EChartsReactProps>(
  () => import("echarts-for-react"),
  { ssr: false }
);

type LineSeries = {
  name: string;
  data: number[];
  color?: string;
};

type LineChartProps = {
  title?: string;
  xData: string[];
  series: LineSeries[];
  height?: number | string;
  minHeight?: number | string;
  customOptions?: EChartsOption;
};

const LineChart: FC<LineChartProps> = ({
  title,
  xData,
  series,
  height = "100%",
  minHeight = "500px",
  customOptions = {}
}) => {
  const options: EChartsOption = {
    title: {
      text: title,
      left: "center"
    },
    tooltip: {
      trigger: "axis"
    },
    xAxis: {
      type: "category",
      data: xData,
      axisLine: {
        lineStyle: {
          color: "gray"
        }
      },
      axisLabel: {
        color: "white"
      }
    },
    yAxis: {
      type: "value",
      axisLine: {
        lineStyle: {
          color: "gray"
        }
      },
      axisLabel: {
        color: "white"
      },
      splitLine: {
        lineStyle: {
          color: "gray"
        }
      }
    },
    series: series.map((s) => ({
      name: s.name,
      type: "line",
      data: s.data,
      smooth: true,
      lineStyle: { width: 3 },
      itemStyle: { color: s.color }
    })),
    ...customOptions
  };

  return (
    <div className="w-full h-full">
      <ReactECharts
        option={options}
        style={{ width: "100%", height, minHeight }}
        opts={{ renderer: "svg" }}
      />
    </div>
  );
};

export default LineChart;
