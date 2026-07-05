import {Button, Card, Col, DatePicker, Form, Input, Row, Select, Statistic} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {playerPros} from "@/utils/utils";
import * as echarts from 'echarts';
import "../../global.less"

const {RangePicker} = DatePicker
const {Option} = Select

@connect(
  state => ({
    ...state.global,
    loading: state.loading.effects["global/fetchPlayerList"]
  })
)
class Player extends React.Component {
  formRef = React.createRef();

  componentDidMount() {
    window.addEventListener('resize', this.resizeCharts)
    this.query().then()
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resizeCharts)
    clearTimeout(this.timelineZoomTimer)
    if (this.classCompareChart) {
      echarts.dispose(this.classCompareChart)
      this.classCompareChart = null
    }
    if (this.timeline) {
      echarts.dispose(this.timeline)
      this.timeline = null
    }
  }

  resizeCharts = () => {
    if (this.classCompareChart) {
      this.classCompareChart.resize()
    }
    if (this.timeline) {
      this.timeline.resize()
    }
  }

  query = async () => {
    const {dispatch} = this.props
    const fieldValue = this.formRef.current.getFieldValue();
    let st, et
    if (fieldValue.time) {
      st = fieldValue.time[0].format("YYYY-MM-DD HH:mm:ss")
      et = fieldValue.time[1].format("YYYY-MM-DD HH:mm:ss")
    }

    await dispatch({
      type: 'global/fetchTimeline',
      payload: {
        st, et,
      }
    });
    this.initTimeline()

    await dispatch({
      type: 'global/fetchPlayerList',
      payload: {
        st, et,
        name: fieldValue.name && fieldValue.name.trim(),
        type: fieldValue.type,
        class: fieldValue.class
      }
    });
    this.initClassCompareChart();
  }

  queryPlayerStats = async (st, et) => {
    const {dispatch} = this.props
    const fieldValue = this.formRef.current.getFieldValue();

    await dispatch({
      type: 'global/fetchPlayerList',
      payload: {
        st, et,
        name: fieldValue.name && fieldValue.name.trim(),
        type: fieldValue.type,
        class: fieldValue.class
      }
    });
    this.initClassCompareChart();
  }

  onReset = async () => {
    this.formRef.current.resetFields();
    this.query().then()
  };

  getClassData(type) {
    const {playerList} = this.props
    let class2num = {};
    playerList && playerList.forEach(v => {
      if (v.type !== type) {
        return
      }
      if (class2num[v.class]) {
        class2num[v.class] = class2num[v.class] + 1;
      } else {
        class2num[v.class] = 1;
      }
    });
    let result = []
    Object.keys(class2num).forEach(key => {
      result.push({
        name: playerPros[key].name + ": " + class2num[key],
        value: class2num[key]
      });
    })
    return result
  }

  getClassCompareData() {
    const {playerList} = this.props
    const classList = playerPros.slice(1)
    const angelData = classList.map(() => 0)
    const demonData = classList.map(() => 0)
    playerList && playerList.forEach(value => {
      const index = classList.findIndex(item => item.class === value.class)
      if (index < 0) {
        return
      }
      if (value.type === 1) {
        angelData[index] += 1
      } else if (value.type === 2) {
        demonData[index] += 1
      }
    })
    return {
      classNames: classList.map(value => value.name),
      angelData,
      demonData,
    }
  }

  initClassCompareChart() {
    try {
      if (this.classCompareChart) {
        echarts.dispose(this.classCompareChart)
      }
      const chartEl = document.getElementById("classCompareChart")
      if (!chartEl) {
        return
      }
      this.classCompareChart = echarts.init(chartEl)
    } catch (e) {
      console.log(e)
      return
    }
    const {classNames, angelData, demonData} = this.getClassCompareData()
    const hasData = angelData.some(Boolean) || demonData.some(Boolean)
    const option = {
      animationDuration: 360,
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        },
        formatter: function (params) {
          const name = params[0] && params[0].axisValue
          const angel = params.find(item => item.seriesName === "天族")
          const demon = params.find(item => item.seriesName === "魔族")
          return `${name}<br/>天族：${Math.abs(angel && angel.value || 0)}<br/>魔族：${Math.abs(demon && demon.value || 0)}`
        }
      },
      legend: {
        data: ['天族', '魔族']
      },
      grid: {
        containLabel: true,
        left: 8,
        right: 30,
        top: 36,
        bottom: 12
      },
      xAxis: {
        type: 'value',
        axisLabel: {
          formatter: function (value) {
            return Math.abs(value)
          }
        }
      },
      yAxis: {
        type: 'category',
        inverse: true,
        data: classNames,
        axisTick: {
          show: false
        }
      },
      graphic: hasData ? [] : [{
        type: 'text',
        left: 'center',
        top: 'middle',
        style: {
          text: '暂无数据',
          fill: '#999',
          fontSize: 14
        }
      }],
      series: [
        {
          name: '天族',
          type: 'bar',
          stack: 'total',
          data: angelData.map(value => -value),
          itemStyle: {
            color: '#52c41a'
          },
          label: {
            show: true,
            position: 'left',
            formatter: function (params) {
              return Math.abs(params.value)
            }
          }
        },
        {
          name: '魔族',
          type: 'bar',
          stack: 'total',
          data: demonData,
          itemStyle: {
            color: '#1890ff'
          },
          label: {
            show: true,
            position: 'right'
          }
        },
      ]
    }
    this.classCompareChart.setOption(option)
  }

  initTimeline() {
    try {
      if (this.timeline) {
        echarts.dispose(this.timeline)
      }
      const chartEl = document.getElementById("timeline")
      if (!chartEl) {
        return
      }
      this.timeline = echarts.init(chartEl)
    }catch (e) {
      console.log(e)
      return
    }
    const {timeline} = this.props
    const timeData = timeline.timeData || []
    const killValue = timeline.killValue || []
    const killedValue = timeline.killedValue || []
    const hasData = timeData.length > 0
    const option = {
      animationDuration: 360,
      grid: {
        containLabel: true,
        left: 40,
        right: 20,
        top: 42,
        bottom: 48
      },
      toolbox: {
        feature: {
          dataZoom: {
            yAxisIndex: 'none'
          },
        }
      },
      legend: {
        data: ['天族击杀数', '魔族击杀数']
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        },
        formatter: function (params) {
          const time = params[0] && params[0].axisValue
          const angel = params.find(item => item.seriesName === "天族击杀数")
          const demon = params.find(item => item.seriesName === "魔族击杀数")
          const angelValue = angel && angel.value || 0
          const demonValue = demon && demon.value || 0
          return `${time}<br/>天族击杀：${angelValue}<br/>魔族击杀：${demonValue}<br/>战斗热度：${angelValue + demonValue}`
        }
      },
      xAxis: {
        show: true,
        type: 'category',
        boundaryGap: true,
        data: timeData,
        axisLabel: {
          formatter: function (value) {
            return moment(value).format("HH:mm")
          }
        },
      },
      yAxis: {
        type: 'value',
        name: '击杀数',
      },
      dataZoom: [
        {
          type: 'slider',
          xAxisIndex: 0,
          height: 18,
          bottom: 12,
        }
      ],
      graphic: hasData ? [] : [{
        type: 'text',
        left: 'center',
        top: 'middle',
        style: {
          text: '暂无战斗趋势',
          fill: '#999',
          fontSize: 14
        }
      }],
      series: [
        {
          name: "天族击杀数",
          type: 'bar',
          barMaxWidth: 12,
          data: killValue,
          itemStyle: {
            color: '#52c41a'
          },
        },
        {
          name: '魔族击杀数',
          type: 'bar',
          barMaxWidth: 12,
          data: killedValue,
          itemStyle: {
            color: '#1890ff'
          },
        },
      ],
    }
    this.timeline.on('datazoom', (params) => {
      clearTimeout(this.timelineZoomTimer)
      this.timelineZoomTimer = setTimeout(() => this.handleTimelineZoom(params), 220)
    })
    this.timeline.on('restore', async () => {
      await this.formRef.current.resetFields(['time'])
      this.queryPlayerStats().then()
    })
    this.timeline.setOption(option)
  }

  handleTimelineZoom = async (params) => {
    const {timeline} = this.props
    const zoom = params.batch && params.batch[0] || params
    const timeData = timeline.timeData || []
    if (!timeData.length) {
      return
    }

    let start = zoom && zoom.startValue
    let end = zoom && zoom.endValue
    if (start === undefined || end === undefined) {
      start = Math.floor((zoom.start || 0) * (timeData.length - 1) / 100)
      end = Math.ceil((zoom.end || 100) * (timeData.length - 1) / 100)
    }

    const startIndex = typeof start === 'number' ? start : timeData.indexOf(start)
    const endIndex = typeof end === 'number' ? end : timeData.indexOf(end)
    const isReset = (zoom.start === 0 && zoom.end === 100) ||
      (startIndex === 0 && endIndex >= timeData.length - 1)
    if (isReset) {
      await this.formRef.current.resetFields(['time'])
      this.queryPlayerStats().then()
      return
    }

    const startTime = typeof start === 'number' ? timeData[start] : start
    const endTime = typeof end === 'number' ? timeData[end] : end
    if (startTime && endTime) {
      await this.formRef.current.setFieldsValue({time: [moment(startTime), moment(endTime)]})
      this.queryPlayerStats(startTime, endTime).then()
    }
  }

  getStatData(data) {
    let angel = 0;
    let demon = 0;
    let other = 0;
    data && data.forEach(v => {
      switch (v.type) {
        case 0:
          other++
          break
        case 1:
          angel++
          break
        case 2:
          demon++
          break
      }
    })
    return {angel, demon, other}
  }

  searchForm() {
    const onFinish = async () => {
      this.query().then()
    };
    const dateFormat = 'YYYY-MM-DD HH:mm:ss';
    return (
      <Form
        layout="inline"
        onFinish={onFinish}
        autoComplete="false"
        ref={this.formRef}
      >
        <Form.Item label="时间" name="time" className="battle-form-range">
          <RangePicker
            format={dateFormat}
            ranges={{
              今天: [moment().startOf('day'), moment().endOf('day')],
              昨天: [moment().subtract(1, 'day').startOf('day'), moment().subtract(1, 'day').endOf('day')],
              前天: [moment().subtract(2, 'day').startOf('day'), moment().subtract(2, 'day').endOf('day')],
              最近3天: [moment().subtract(2, 'day').startOf('day'), moment().endOf('day')],
              最近7天: [moment().subtract(6, 'day').startOf('day'), moment().endOf('day')],
            }}
            allowClear
            showTime={{defaultValue: moment('00:00:00', 'HH:mm:ss')}}
            onChange={(d, ds) => this.query(d, ds)}
            style={{width: '100%'}}
          />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            搜索
          </Button>
          &nbsp;&nbsp;
          <Button onClick={this.onReset}>
            重置
          </Button>
        </Form.Item>
      </Form>)
  }

  render() {
    const {playerList} = this.props
    const players = playerList || []
    const statData = this.getStatData(playerList)
    const knownPlayers = statData.angel + statData.demon
    const angelRate = knownPlayers ? `${((statData.angel / knownPlayers) * 100).toFixed(1)}%` : '0%'
    const demonRate = knownPlayers ? `${((statData.demon / knownPlayers) * 100).toFixed(1)}%` : '0%'
    return (
      <PageContainer title={false}>
        <Card className="battle-toolbar">
          {this.searchForm()}
        </Card>
        <Row gutter={[12, 12]}>
          <Col xs={24} xl={12}>
            <Card title="玩家统计" className="battle-section-card battle-overview-stats">
              <Row gutter={[12, 18]}>
                <Col xs={24} md={8} xl={24}>
                  <Statistic title="玩家总数" value={players.length}/>
                </Col>
                <Col xs={12} md={8} xl={12}>
                  <Statistic title={`天族 ${angelRate}`} value={statData.angel} valueStyle={{color: "#52c41a"}}/>
                </Col>
                <Col xs={12} md={8} xl={12}>
                  <Statistic title={`魔族 ${demonRate}`} value={statData.demon} valueStyle={{color: "#1890ff"}}/>
                </Col>
                <Col xs={24} md={8} xl={24}>
                  <Statistic title="未识别/其它" value={statData.other} valueStyle={{color: "#fa8c16"}}/>
                </Col>
              </Row>
            </Card>
          </Col>
          <Col xs={24} xl={12}>
            <Card title="天魔职业对比" className="battle-section-card">
              <div id="classCompareChart" className="battle-chart-compare"/>
            </Card>
          </Col>
        </Row>
        <Card title="战斗趋势" className="battle-section-card battle-primary-card">
          <div id="timeline" className="battle-chart-main"/>
        </Card>
      </PageContainer>
    );
  }
}

export default Player;
