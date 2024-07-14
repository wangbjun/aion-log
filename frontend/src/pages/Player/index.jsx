import {Button, Card, Col, DatePicker, Empty, Form, Input, Row, Select, Statistic, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {Link} from 'umi';
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

  constructor(props) {
    super(props);
    this.columns = [
      {
        title: "玩家",
        dataIndex: 'name',
        key: 'name',
        sorter: function (a, b) {
          return a.name.localeCompare(b.name)
        },
        render: this.renderName,
        width: "25%",
      },
      {
        title: "种族",
        dataIndex: 'type',
        key: 'type',
        width: "10%",
        sorter: function (a, b) {
          return a.type - b.type
        },
        render: function (value) {
          if (value === 0) {
            return <Tag color="orange">其它</Tag>
          }
          if (value === 1) {
            return <Tag color="green">天族</Tag>
          }
          if (value === 2) {
            return <Tag color="blue">魔族</Tag>
          }
        }
      },
      {
        title: "职业",
        dataIndex: 'class',
        key: 'class',
        width: "8%",
        sorter: function (a, b) {
          return a.class - b.class
        },
        render: function (value) {
          return <img src={require("../../assets/" + playerPros[value].logo)} width={30}/>
        }
      },
      {
        title: "技能次数",
        dataIndex: 'skill_count',
        key: 'skill_count',
        width: "8%",
        sorter: function (a, b) {
          return a.skill_count - b.skill_count
        },
      },
      {
        title: "暴击率",
        dataIndex: 'critical_ratio',
        key: 'critical_ratio',
        width: "8%",
        sorter: function (a, b) {
          return a.critical_ratio - b.critical_ratio
        },
        render: function (value) {
          return (value*100).toFixed(1)+"%"
        }
      },
      {
        title: "击杀数",
        dataIndex: 'kill_count',
        key: 'kill_count',
        width: "8%",
        sorter: function (a, b) {
          return a.kill_count - b.kill_count
        },
      },
      {
        title: "死亡数",
        dataIndex: 'death_count',
        key: 'death_count',
        width: "8%",
        sorter: function (a, b) {
          return a.death_count - b.death_count
        },
      },
      {
        title: "最后更新时间",
        dataIndex: 'time',
        key: 'time',
        width: "20%",
        sorter: function (a, b) {
          return moment(a.time).isAfter(moment(b.time))
        },
        render: function (value) {
          return moment(value).format("YYYY-MM-DD HH:mm:ss")
        }
      },
    ];
  }

  renderName = (value) => {
    return <Link target="_blank" to={`/log?player=${value}`}>{value}</Link>
  }

  componentDidMount() {
    this.query().then()
  }

  query = async () => {
    const {dispatch} = this.props
    const fieldValue = this.formRef.current.getFieldValue();
    let st = moment().subtract(6, 'day').startOf('day').format("YYYY-MM-DD HH:mm:ss")
    let et = moment().endOf('day').format("YYYY-MM-DD HH:mm:ss")
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
        name: fieldValue.name,
        type: fieldValue.type,
        class: fieldValue.class
      }
    });
    this.initAngelPie();
    this.initDemonPie();
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

  initAngelPie() {
    try {
      if (!this.angelPie) {
        this.angelPie = echarts.init(document.getElementById("angelPie"))
      }
    }catch (e) {
      console.log(e)
      return
    }
    const option = {
      title: {
        text: "天族职业分布"
      },
      tooltip: {
        trigger: 'item',
        formatter: '{b0}'
      },
      legend: {
        orient: 'vertical',
        left: 'right',
        align: 'left'
      },
      series: [
        {
          type: 'pie',
          radius: '80%',
          data: this.getClassData(1),
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          },
          label: {
            show: true,
            position: 'inside',
            formatter: '{d}%'
          }
        }
      ]
    }
    this.angelPie.setOption(option)
  }

  initDemonPie() {
    try {
      if (!this.demonPie) {
        this.demonPie = echarts.init(document.getElementById("demonPie"))
      }
    }catch (e) {
      console.log(e)
      return
    }
    const option = {
      title: {
        text: "魔族职业分布"
      },
      tooltip: {
        trigger: 'item',
        formatter: '{b0}'
      },
      legend: {
        orient: 'vertical',
        left: 'right',
        align: 'left'
      },
      series: [
        {
          type: 'pie',
          radius: '80%',
          data: this.getClassData(2),
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          },
          label: {
            show: true,
            position: 'inside',
            formatter: '{d}%'
          }
        }
      ]
    }
    this.demonPie.setOption(option)
  }

  initTimeline() {
    try {
      if (this.timeline) {
        echarts.dispose(this.timeline)
      }
      this.timeline = echarts.init(document.getElementById("timeline"))
    }catch (e) {
      console.log(e)
      return
    }
    const {timeline} = this.props
    const option = {
      grid: {
        left: 30,
        right: 20,
        top: '10%',
        bottom: 5
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
      },
      xAxis: {
        show: false,
        type: 'category',
        boundaryGap: true,
        data: timeline.timeData,
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: "天族击杀数",
          type: 'line',
          data: timeline.killValue,
        },
        {
          name: '魔族击杀数',
          type: 'line',
          data: timeline.killedValue,
        },
      ],
    }
    this.timeline.on('datazoom', async (params) => {
      let start = params.batch && params.batch[0].startValue
      let end = params.batch && params.batch[0].endValue
      let startTime = timeline.timeData && timeline.timeData[start]
      let endTime = timeline.timeData && timeline.timeData[end]
      if (startTime && endTime) {
        this.formRef.current.setFieldsValue({time: [moment(startTime), moment(endTime)]})
        this.query().then()
      }
    })
    this.timeline.setOption(option)
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
        <Form.Item label="时间" name="time">
          <RangePicker
            format={dateFormat}
            defaultValue={[moment().subtract(6, 'day').startOf('day'), moment().endOf('day')]}  // 设置默认值为最近7天
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
            style={{width: 350}}
          />
        </Form.Item>
        <Form.Item label="玩家" name="name">
          <Input allowClear placeholder="请输入" style={{width: 150}}/>
        </Form.Item>
        <Form.Item label="种族" name="type">
          <Select
            allowClear
            showSearch
            style={{width: 100}}
            placeholder="请选择种族"
            optionFilterProp="children"
            filterOption={(input, option) =>
              option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
            }
            onSelect={() => this.query()}
          >
            <Option value="1">天族</Option>
            <Option value="2">魔族</Option>
            <Option value="0">其它</Option>
          </Select>
        </Form.Item>
        <Form.Item label="职业" name="class">
          <Select
            allowClear
            showSearch
            style={{width: 100}}
            placeholder="请选择职业"
            optionFilterProp="children"
            filterOption={(input, option) =>
              option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
            }
            onSelect={() => this.query()}
          >
            {playerPros.map((v, k) =>
              <Option value={k} key={k}>{v.name}</Option>
            )}
          </Select>
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            搜索
          </Button>
          &nbsp;&nbsp;
          <Button type="primary" onClick={this.onReset}>
            重置
          </Button>
        </Form.Item>
      </Form>)
  }

  render() {
    const {playerList, loading} = this.props
    const statData = this.getStatData(playerList)
    return (
      <PageContainer>
        <Card extra={this.searchForm()} >
          <Row>
            <Col span={7}>
              <Card title="种族">
                <Row gutter={16}>
                  <Col span={6}>
                    <Statistic title="天魔总数" value={statData.angel + statData.demon} valueStyle={{color: "red"}}/>
                  </Col>
                  <Col span={6}>
                    <Statistic title="天族" value={statData.angel} valueStyle={{color: "green"}}/>
                  </Col>
                  <Col span={6}>
                    <Statistic title="魔族" value={statData.demon} valueStyle={{color: "blue"}}/>
                  </Col>
                  <Col span={6}>
                    <Statistic title="其它" value={statData.other} valueStyle={{color: "grey"}}/>
                  </Col>
                </Row>
              </Card>
              <Card>
                <div id="angelPie" style={{height: '250px'}}/>
              </Card>
              <Card>
                <div id="demonPie" style={{height: '250px'}}/>
              </Card>
            </Col>
            <Col span={17}>
              <Card>
                <div id="timeline" style={{height: '250px'}}/>
              </Card>
              <Table
                bordered
                size="small"
                columns={this.columns}
                dataSource={playerList}
                rowKey={(record) => {
                  return record.id
                }}
                pagination={{
                  defaultPageSize: 15,
                  pageSizeOptions: ['50', '100', '200', '500'],
                  showTotal: (total) => `共${total}条记录`,
                }}
                loading={loading}
              />
            </Col>
          </Row>
        </Card>
      </PageContainer>
    );
  }
}

export default Player;
