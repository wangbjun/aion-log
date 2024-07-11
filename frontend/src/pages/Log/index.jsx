import {Button, Card, DatePicker, Form, Input, Select, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {parse} from 'querystring'
import {getTypeColor, playerPros} from "@/utils/utils";
import "../../global.less"

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
    pageSize: 50,
    valueGe: "",
    valueLe: ""
  }

  constructor(props) {
    super(props);
    this.columns = [
      {
        title: "时间",
        dataIndex: 'time',
        key: 'time',
        width: 145,
        render: (value, row) => {
          return moment(value).format("YYYY-MM-DD HH:mm:ss")
        }
      },
      {
        title: "原始日志",
        dataIndex: 'raw_msg',
        key: 'raw_msg',
        render: (value, row) => {
          if (row.skill === "kill" || row.skill === "killed") {
            return <div style={{color: "deeppink"}}>{value}</div>;
          }
          if (!row.skill) {
            return <div>{value}</div>;
          }
          const [color,typeName] = getTypeColor(row.player_type)
          let results = []
          const parts = value.split(row.player);
          results.push(parts[0])
          results.push((<span key={row.id+row.player}><Tag className="custom-tag" color={color}>{typeName}</Tag><Tag className="custom-tag">{playerPros[row.player_class].name}</Tag><a onClick={() => this.searchPlayer(row.player)}>{row.player}</a></span>))

          const parts2 = parts[1].split(row.skill);
          results.push(parts2[0]);
          if (row.skill !== "attack") {
            results.push(<span style={{color: "red", fontWeight: "bold"}} key={1}>{row.skill}</span>)
          }
          if (row.target !== "" && parts2[1]) {
            const parts3 = parts2[1].split(row.target);
            const [color,typeName] = getTypeColor(row.target_type)
            results.push(parts3[0]);
            results.push((<span key={row.id+row.target}><Tag className="custom-tag" color={color}>{typeName}</Tag><Tag className="custom-tag">{playerPros[row.target_class].name}</Tag><a onClick={() => this.searchPlayer(row.target)}>{row.target}</a></span>))
            results.push(parts3[1]);
          }else {
            results.push(parts2[1])
          }
          return <div>{results}</div>;
        },
      },
      {
        title: "数值",
        dataIndex: 'value',
        key: 'value',
      },
    ];
  }

  componentDidMount() {
    const parsedUrlQuery = parse(window.location.href.split('?')[1]);
    let param = parsedUrlQuery.player
    if (param) {
      if (param.endsWith("#/")) {
        param = param.substring(0, param.lastIndexOf("#/"))
      }
      this.formRef.current.setFieldsValue({player: param, target: param})
    }
    this.query().then()
  }

  async searchPlayer(player) {
    await this.formRef.current.setFieldsValue({player: player, target: player})
    await this.setState({page: 1})
    this.props.history.push("/log?player=" + player)
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
    const {page, pageSize, valueGe, valueLe} = this.state
    console.log(valueGe, valueLe)
    dispatch({
      type: 'global/fetchLogList',
      payload: {
        page,
        pageSize,
        st: ds && ds[0] || st,
        et: ds && ds[1] || et,
        player, target, skill,
        value: valueGe||valueLe ? valueGe+"-"+valueLe : "",
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

  inputValue = async (e, type) => {
    if (type === "a") {
      await this.setState({valueGe: e.target.value})
    } else if (type === "b") {
      await this.setState({valueLe: e.target.value})
    }
  }

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
              前天: [moment().subtract(2, 'day').startOf('day'), moment().subtract(2, 'day').endOf('day')],
              最近3天: [moment().subtract(2, 'day').startOf('day'), moment().endOf('day')],
              最近7天: [moment().subtract(6, 'day').startOf('day'), moment().endOf('day')],
            }}
            allowClear
            showTime={{defaultValue: moment('00:00:00', 'HH:mm:ss')}}
            onChange={(d, ds) => this.query(d, ds)}
            style={{width: 300}}
          />
        </Form.Item>
        <Form.Item label="技能" name="skill">
          <Input allowClear placeholder="请输入" style={{width: 150}}/>
        </Form.Item>
        <Form.Item label="玩家" name="player">
          <Input allowClear placeholder="请输入" style={{width: 150}}/>
        </Form.Item>
        <Form.Item label="对象" name="target">
          <Input allowClear placeholder="请输入" style={{width: 150}}/>
        </Form.Item>
        <Form.Item label="数值" name="value">
          <Input allowClear placeholder=">=" style={{width: 70}} onChange={(e)=>this.inputValue(e, "a")}/> - <Input allowClear placeholder="<=" style={{width: 70}} onChange={(e)=>this.inputValue(e, "b")}/>
        </Form.Item>
        <Form.Item label="排序" name="sort">
          <Select
            allowClear
            showSearch
            placeholder="请选择"
            optionFilterProp="children"
            filterOption={(input, option) =>
              option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
            }
            onSelect={() => this.query()}
            style={{width: 100}}
          >
            <Option value="time">时间</Option>
            <Option value="value">数值</Option>
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
      defaultPageSize: 50,
      total: logList.total,
      pageSizeOptions: ['50', '100', '200', '500'],
      showTotal: (total) => `共${total}条记录`,
      onChange: async (page, pageSize) => {
        await this.setState({page, pageSize})
        this.query().then()
      }
    }
    let color = "row-odd"
    const rowClassName = (record, index) => {
      if (index === 0) {
        return color
      }
      if (record.time === logList.list[index - 1].time && record.player === logList.list[index - 1].player) {
        return color
      } else {
        if (color === "row-odd") {
          color = "row-even"
        } else {
          color = "row-odd"
        }
      }
      return color
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
            rowClassName={rowClassName}
          />
        </Card>
      </PageContainer>
    );
  }
}

export default Log;
