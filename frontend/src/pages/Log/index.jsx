import {Button, Card, DatePicker, Form, Input, Select, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {parse} from 'querystring'
import {playerPros} from "@/utils/utils";

const {RangePicker} = DatePicker
const {Option} = Select

@connect(
  state => ({
    ...state.global,
    loading: state.loading.effects["global/fetchLogList"]
  })
)
class Log extends React.Component {
  formRef = React.createRef();

  state = {
    page: 1,
    pageSize: 20,
  }

  constructor(props) {
    super(props);
    this.columns = [
      {
        title: "时间",
        dataIndex: 'time',
        key: 'time',
        width: "10%",
        render: this.renderTime
      },
      {
        title: "玩家",
        dataIndex: 'player',
        key: 'player',
        width: "18%",
        render: function (value, row) {
          let color = "grey"
          let typeName = ""
          if (row.player_type === 1) {
            color = "green"
            typeName = "天族"
          } else if (row.player_type === 2) {
            color = "blue"
            typeName = "魔族"
          } else if (row.player_type === 0) {
            color = "orange"
            typeName = "其它"
          }
          return <div><Tag className="custom-tag" color={color}>{typeName}</Tag><Tag className="custom-tag">{playerPros[row.player_class].name}</Tag><span>{value}</span></div>
        }
      },
      {
        title: "对象",
        dataIndex: 'target',
        key: 'target',
        width: "18%",
        render: function (value, row) {
          if (!value) {
            return <div>{value}</div>;
          }
          let color = "grey"
          let typeName = ""
          if (row.target_type === 1) {
            color = "green"
            typeName = "天族"
          } else if (row.target_type === 2) {
            color = "blue"
            typeName = "魔族"
          } else if (row.target_type === 0) {
            color = "orange"
            typeName = "其它"
          }
          return <div><Tag className="custom-tag" color={color}>{typeName}</Tag><Tag className="custom-tag">{playerPros[row.target_class].name}</Tag><span>{value}</span></div>
        }
      },
      {
        title: "伤害",
        dataIndex: 'value',
        key: 'value',
        width: "6%",
      },
      {
        title: "原始日志",
        dataIndex: 'raw_msg',
        key: 'raw_msg',
        width: "50%",
        render: function (value, row) {
          if (value.indexOf("打倒了。") !== -1 || value.indexOf("攻击而终结。") !== -1) {
            return <div style={{color: "deeppink"}}>{value}</div>;
          }
          if (!row.skill) {
            return <div>{value}</div>;
          }
          let results = []
          const parts = value.split(row.skill);
          results.push(parts[0])
          if (row.skill !== "普通攻击") {
            results.push(<span style={{color: "red", fontWeight: "bold"}} key={1}>{row.skill}</span>)
          }
          results.push(parts[1])
          return <div>{results}</div>;
        }
      },
    ];
  }

  renderTime = (value, row) => {
    const time = moment(value).format("YYYY-MM-DD HH:mm:ss")
    return (<a onClick={() => this.searchTime(row)}>{time}</a>)
  }

  componentDidMount() {
    const parsedUrlQuery = parse(window.location.href.split('?')[1]);
    let param = parsedUrlQuery.player
    if (param) {
      if (param.endsWith("#/")) {
        param = param.substring(0, param.lastIndexOf("#/"))
      }
      this.formRef.current.setFieldsValue({player: param})
    }
    this.query().then()
  }

  async searchTime(row) {
    await this.formRef.current.setFieldsValue({time: [moment(row.time), moment(row.time)], player: row.player})
    await this.setState({page: 1})
    this.props.history.push("/log?player=" + row.player)
    this.query().then()
  }

  query = async (d, ds) => {
    const {dispatch} = this.props
    const fieldValue = this.formRef.current.getFieldValue();
    let st = fieldValue.time && fieldValue.time[0].format("YYYY-MM-DD HH:mm:ss")
    let et = fieldValue.time && fieldValue.time[1].format("YYYY-MM-DD HH:mm:ss")
    let player = fieldValue.player && fieldValue.player.trim()
    let target = fieldValue.target && fieldValue.target.trim()
    let skill = fieldValue.skill && fieldValue.skill.trim()
    let value = fieldValue.value && fieldValue.value.trim()
    const {page, pageSize} = this.state
    dispatch({
      type: 'global/fetchLogList',
      payload: {
        page,
        pageSize,
        st: ds && ds[0] || st,
        et: ds && ds[1] || et,
        player, target, skill,value,
        sort: fieldValue.sort
      },
    });
  }

  onReset = async () => {
    await this.formRef.current.resetFields();
    await this.setState({page: 1})
    this.props.history.push("/log")
    this.query().then()
  };

  searchForm() {
    const dateFormat = 'YYYY-MM-DD HH:mm:ss';
    const onFinish = async () => {
      await this.setState({page: 1})
      this.query().then()
    };
    return (
      <Form
        layout="inline"
        onFinish={onFinish}
        autoComplete="false"
        ref={this.formRef}
        style={{overflow: "right"}}
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
        <Form.Item label="技能" name="skill">
          <Input allowClear placeholder="请输入" style={{ width: 150 }} />
        </Form.Item>
        <Form.Item label="玩家" name="player">
          <Input allowClear placeholder="请输入" style={{ width: 150 }}/>
        </Form.Item>
        <Form.Item label="对象" name="target">
          <Input allowClear placeholder="请输入" style={{ width: 150 }}/>
        </Form.Item>
        <Form.Item label="伤害大于" name="value">
          <Input allowClear placeholder="请输入" style={{ width: 100 }}/>
        </Form.Item>
        <Form.Item label="排序" name="sort">
          <Select
            allowClear
            showSearch
            placeholder="请选择排序"
            optionFilterProp="children"
            filterOption={(input, option) =>
              option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
            }
            onSelect={() => this.query()}
            style={{ width: 100 }}
          >
            <Option value="time">时间</Option>
            <Option value="value">伤害</Option>
            <Option value="skill">技能</Option>
            <Option value="player">玩家</Option>
            <Option value="target">对象</Option>
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
    const {page, pageSize} = this.state
    const {logList, loading} = this.props
    const pagination = {
      current: page,
      pageSize: pageSize,
      defaultPageSize: 20,
      total: logList.total,
      pageSizeOptions:['50', '100', '200', '500'],
      showTotal: (total) => `共${total}条记录`,
      onChange: async (page, pageSize) => {
        await this.setState({page, pageSize})
        this.query().then()
      },
      hideOnSinglePage: true
    }
    return (
      <PageContainer>
        <Card extra={this.searchForm()}>
          <Table
            bordered
            size="small"
            columns={this.columns}
            dataSource={logList.list}
            rowKey={(record) => {
              return record.id
            }}
            pagination={pagination}
            loading={loading}
          />
        </Card>
      </PageContainer>
    );
  }
}

export default Log;
