import {Button, Card, DatePicker, Form, Input, Select, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {parse} from 'querystring'

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
    pageSize: 15,
  }

  constructor(props) {
    super(props);
    this.columns = [
      {
        title: "时间",
        dataIndex: 'time',
        key: 'time',
        render: this.renderTime
      },
      {
        title: "玩家",
        dataIndex: 'player',
        key: 'player',
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
          return <div><Tag className="custom-tag" color={color}>{typeName}</Tag><span>{value}</span></div>
        }
      },
      {
        title: "被玩家",
        dataIndex: 'target_player',
        key: 'target_player',
        render: function (value, row) {
          let color = "grey"
          let typeName = ""
          if (row.target_player_type === 1) {
            color = "green"
            typeName = "天族"
          } else if (row.target_player_type === 2) {
            color = "blue"
            typeName = "魔族"
          } else if (row.target_player_type === 0) {
            color = "orange"
            typeName = "其它"
          }
          return <span><Tag className="custom-tag" color={color}>{typeName}</Tag>{value}</span>
        }
      },
      {
        title: "伤害",
        dataIndex: 'damage',
        key: 'damage'
      },
      {
        title: "原始日志",
        dataIndex: 'origin_desc',
        key: 'origin_desc',
        width: "50%",
        render: function (value, row) {
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
    } else {
      this.formRef.current.setFieldsValue({
        time: [moment().subtract(6, 'day').startOf('day'), moment().endOf('day')]
      })
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
    let skill = fieldValue.skill && fieldValue.skill.trim()
    const {page, pageSize} = this.state
    dispatch({
      type: 'global/fetchLogList',
      payload: {
        page,
        pageSize,
        st: ds && ds[0] || st,
        et: ds && ds[1] || et,
        player, skill,
        sort: fieldValue.sort
      },
    });
  }

  onReset = async () => {
    await this.formRef.current.resetFields();
    await this.formRef.current.setFieldsValue({
      time: [moment().subtract(6, 'day').startOf('day'), moment().endOf('day')]
    })
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
        <Form.Item label="时间" name="time" style={{marginTop: "5px"}}>
          <RangePicker
            format={dateFormat}
            ranges={{
              今天: [moment().startOf('day'), moment().endOf('day')],
              昨天: [moment().subtract(1, 'day').startOf('day'), moment().subtract(1, 'day').endOf('day')],
              最近三天: [moment().subtract(2, 'day').startOf('day'), moment().endOf('day')],
              最近一周: [moment().subtract(6, 'day').startOf('day'), moment().endOf('day')],
            }}
            allowClear
            showTime={{defaultValue: moment('00:00:00', 'HH:mm:ss')}}
            onChange={(d, ds) => this.query(d, ds)}
          />
        </Form.Item>
        <Form.Item label="技能" name="skill" style={{marginTop: "5px"}}>
          <Input allowClear placeholder="请输入"/>
        </Form.Item>
        <Form.Item label="玩家" name="player" style={{marginTop: "5px"}}>
          <Input allowClear placeholder="请输入"/>
        </Form.Item>
        <Form.Item label="排序" name="sort" style={{marginTop: "5px"}}>
          <Select
            allowClear
            showSearch
            placeholder="请选择排序"
            optionFilterProp="children"
            filterOption={(input, option) =>
              option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
            }
            onSelect={() => this.query()}
          >
            <Option value="time">时间</Option>
            <Option value="damage">伤害</Option>
            <Option value="skill">技能</Option>
            <Option value="player">玩家</Option>
            <Option value="target_player">被玩家</Option>
          </Select>
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" style={{marginTop: "5px"}}>
            搜索
          </Button>
          &nbsp;&nbsp;
          <Button type="primary" onClick={this.onReset} style={{marginTop: "5px"}}>
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
      defaultPageSize: 15,
      total: logList.total,
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
