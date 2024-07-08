import {Button, Card, Col, DatePicker, Empty, Form, Input, message, Row, Select, Statistic, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {Link} from 'umi';
import {playerPros} from "@/utils/utils";
import { Pie } from '@ant-design/plots';

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
        width: 300,
      },
      {
        title: "种族",
        dataIndex: 'type',
        key: 'type',
        width: 100,
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
        width: 100,
        sorter: function (a, b) {
          return a.class - b.class
        },
        render: function (value) {
          return <img src={require("../../assets/" + playerPros[value].logo)} width={30}/>
        }
      },
      {
        title: "攻击次数",
        dataIndex: 'skill_count',
        key: 'skill_count',
        width: 100,
        sorter: function (a, b) {
          return a.skill_count - b.skill_count
        },
      },
      {
        title: "击杀数",
        dataIndex: 'kill_count',
        key: 'kill_count',
        width: 100,
        sorter: function (a, b) {
          return a.kill_count - b.kill_count
        },
      },
      {
        title: "死亡数",
        dataIndex: 'death_count',
        key: 'death_count',
        width: 100,
        sorter: function (a, b) {
          return a.death_count - b.death_count
        },
      },
      {
        title: "最后更新时间",
        dataIndex: 'time',
        key: 'time',
        width: 180,
        defaultSortOrder: 'descend',
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
    return <Link to={`/log?player=${value}`}>{value}</Link>
  }

  componentDidMount() {
    this.query()
  }

  query = () => {
    const {dispatch} = this.props
    const fieldValue = this.formRef.current.getFieldValue();
    let st = fieldValue.time && fieldValue.time[0].format("YYYY-MM-DD HH:mm:ss")
    let et = fieldValue.time && fieldValue.time[1].format("YYYY-MM-DD HH:mm:ss")
    dispatch({
      type: 'global/fetchPlayerList',
      payload: {
        st, et,
        name: fieldValue.name,
        type: fieldValue.type,
        class: fieldValue.class
      }
    });
  }

  onReset = () => {
    this.formRef.current.resetFields();
    this.query()
  };

  searchForm() {
    const onFinish = async () => {
      this.query()
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
            ranges={{
              今天: [moment().startOf('day'), moment().endOf('day')],
              昨天: [moment().subtract(1, 'day').startOf('day'), moment().subtract(1, 'day').endOf('day')],
              最近3天: [moment().subtract(2, 'day').startOf('day'), moment().endOf('day')],
              最近7天: [moment().subtract(6, 'day').startOf('day'), moment().endOf('day')],
            }}
            allowClear
            showTime={{defaultValue: moment('00:00:00', 'HH:mm:ss')}}
            onChange={(d, ds) => this.query(d, ds)}
            style={{ width: 300 }}
          />
        </Form.Item>
        <Form.Item label="玩家" name="name">
          <Input allowClear placeholder="请输入" style={{ width: 150 }}/>
        </Form.Item>
        <Form.Item label="种族" name="type" >
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
        <Form.Item label="职业" name="class" >
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

  getStatData(data) {
    let tian = 0;
    let mo = 0;
    let other = 0;
    data.forEach(v => {
      switch (v.type) {
        case 0:
          other++
          break
        case 1:
          tian++
          break
        case 2:
          mo++
          break
      }
    })
    return {tian, mo, other}
  }

  getClassData(data) {
    let class2num = {};
    data.forEach(v => {
      if (v.type === 0) {
        return
      }
      if (class2num[v.class]) {
        class2num[v.class] =  class2num[v.class] + 1;
      } else {
        class2num[v.class] = 1;
      }
    });
    let result = []
    Object.keys(class2num).forEach(key => {
      result.push({
        type: playerPros[key].name + ": " +class2num[key],
        value: class2num[key]
      });
    })
    return result
  }

  getServerData(data) {
    let server2num = new Map();
    data.forEach(v => {
      const parts = v.name.split("-");
      if (parts.length !== 2) {
        return
      }
      let serverName = parts[1]
      if (server2num.has(serverName)) {
        server2num.set(serverName, server2num.get(serverName) + 1);
      } else {
        server2num.set(serverName, 1);
      }
    });
    return server2num
  }

  render() {
    const {playerList, loading} = this.props
    const statData = this.getStatData(playerList)
    const classData = this.getClassData(playerList)
    const serverData = this.getServerData(playerList)
    const config = {
      data: classData,
      autoFit: true,
      angleField: 'value',
      colorField: 'type',
      label: {
        text: 'type',
        position: 'inside',
        formatter: (text, datum, index, data) => {
          return text.split(":")[0]
        }
      },
      legend: {
        show: false,
        color: {
          position: 'right',
          rowPadding: 3,
        }
      },
      tooltip: (
        d, // 每一个数据项
        index, // 索引
        data, // 完整数据
        column, // 通道
      ) => ({
        value: `人数:${d.value},占比:${(d.value/(statData.tian+statData.mo)*100).toFixed(0)}%`,
      })
    };
    return (
      <PageContainer>
        <Card extra={this.searchForm()}>
          <Row>
            <Col span={8}>
              <Card title="种族">
                <Row gutter={24}>
                  <Col span={6}>
                    <Statistic title="总数" value={statData.tian+statData.mo} style={{padding: "12px"}} valueStyle={{color: "red"}}/>
                  </Col>
                  <Col span={6}>
                    <Statistic title="天族" value={statData.tian} style={{padding: "12px"}} valueStyle={{color: "green"}}/>
                  </Col>
                  <Col span={6}>
                    <Statistic title="魔族" value={statData.mo} style={{padding: "12px"}} valueStyle={{color: "blue"}}/>
                  </Col>
                  <Col span={6}>
                    <Statistic title="其它" value={statData.other} style={{padding: "12px"}} valueStyle={{color: "grey"}}/>
                  </Col>
                </Row>
              </Card>
              <Card title="职业">
                <Row gutter={24}>
                  { classData.length ? <Pie {...config} />: <Empty/>}
                </Row>
              </Card>
              <Card title="区服">
                <Row gutter={24}>
                  {
                    Array.from(serverData.entries()).map((v, k) => {
                      return  <Col span={6} key={k}>
                        <Statistic title={v[0]} value={v[1]} style={{padding: "12px"}} valueStyle={{color: "green"}}/>
                      </Col>
                    })
                  }
                </Row>
              </Card>
            </Col>
            <Col span={16}>
              <Table
                bordered
                size="small"
                columns={this.columns}
                dataSource={playerList}
                rowKey={(record) => {
                  return record.id
                }}
                pagination={{
                  defaultPageSize: 20,
                  hideOnSinglePage: true,
                  pageSizeOptions:['50', '100', '200', '500'],
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
