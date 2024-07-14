import {Button, Card, Col, DatePicker, Form, Image, Input, Row, Select, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {parse} from 'querystring'
import {getTypeColor, playerPros} from "@/utils/utils";
import "../../global.less"
import {queryPlayer} from "@/services/api";

const {RangePicker} = DatePicker
const {Option} = Select

@connect(
  state => ({
    ...state.global,
    loading: state.loading.effects["global/fetchLogData"],
    loadingTop: state.loading.effects["global/fetchClassTop"]
  })
)
class Log extends React.Component {
  formRef = React.createRef();

  state = {
    page: 1,
    pageSize: 50,
    valueGe: "",
    valueLe: "",
    queryClass: "1"
  }

  constructor(props) {
    super(props);
    this.columns = [
      {
        title: "时间",
        dataIndex: 'time',
        key: 'time',
        width: "15%",
        render: (value, row) => {
          return moment(value).format("YYYY-MM-DD HH:mm:ss")
        }
      },
      {
        title: "战斗信息",
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
            results.push((<span key={row.id+row.target}><Tag className="custom-tag" color={color}>{typeName}</Tag><Tag className="custom-tag">{playerPros[row.target_class].name}</Tag><a onClick={() => this.searchTarget(row.target)}>{row.target}</a></span>))
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

    this.columnsClassTop = [
      {
        title: "技能",
        dataIndex: 'skill',
        key: 'skill',
        render: (value, row) => {
          return <a onClick={async () => {
            await this.formRef.current.setFieldsValue({skill: value})
            await this.formRef.current.setFieldsValue({sort: "value"})
            this.query().then()
          }}>{value}</a>;
        }
      },
      {
        title: "次数",
        dataIndex: 'count',
        key: 'count',
        defaultSortOrder: "descend",
        sorter: function (a, b) {
          return a.count - b.count
        },
      },
      {
        title: "暴击率",
        dataIndex: 'critical',
        key: 'critical',
        sorter: function (a, b) {
          return a.critical - b.critical
        },
        render: function (value, row) {
          return value && (value * 100).toFixed(1)+"%"
        }
      },
      {
        title: "最高伤害",
        dataIndex: 'damage',
        key: 'damage',
        sorter: function (a, b) {
          return a.damage - b.damage
        },
      },
      {
        title: "平均伤害",
        dataIndex: 'average',
        key: 'average',
        render: (value, row) => {
          return value.toFixed(0)
        },
        sorter: function (a, b) {
          return a.average - b.average
        },
      },
    ]
  }

  componentDidMount() {
    const parsedUrlQuery = parse(window.location.href.split('?')[1]);
    let player = parsedUrlQuery.player
    let target = parsedUrlQuery.target
    if (player) {
      if (player.endsWith("#/")) {
        player = player.substring(0, player.lastIndexOf("#/"))
      }
      this.formRef.current.setFieldsValue({player: player})
    }
    if (target) {
      if (target.endsWith("#/")) {
        target = target.substring(0, target.lastIndexOf("#/"))
      }
      this.formRef.current.setFieldsValue({target: target})
    }
    this.query().then()
  }

  async searchPlayer(player) {
    await this.formRef.current.setFieldsValue({player: player})
    await this.setState({page: 1})
    this.props.history.push("/log?player=" + player)
    this.query().then()
  }

  async searchTarget(player) {
    await this.formRef.current.setFieldsValue({target: player})
    await this.setState({page: 1})
    this.props.history.push("/log?target=" + player)
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
    let banPlayer = fieldValue.banPlayer && fieldValue.banPlayer.join(",")
    const {page, pageSize, valueGe, valueLe, queryClass} = this.state

    await dispatch({
      type: 'global/fetchLogData',
      payload: {
        page,
        pageSize,
        st: ds && ds[0] || st,
        et: ds && ds[1] || et,
        player, target, skill, banPlayer,
        value: valueGe||valueLe ? valueGe+"-"+valueLe : "",
        sort: fieldValue.sort,
      },
    });
    this.queryClassTop(queryClass)
  }

  queryClassTop = (queryClass ) => {
    const {dispatch} = this.props
    this.setState({queryClass})
    dispatch({
      type: 'global/fetchClassTop',
      payload: {
        class: queryClass
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
            style={{width: 350}}
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
            placeholder="排序"
            optionFilterProp="children"
            filterOption={(input, option) =>
              option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
            }
            onSelect={() => this.query()}
            style={{width: 70}}
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
    const {page, pageSize, queryClass} = this.state
    const {logData, loading, classTop, loadingTop} = this.props
    const pagination = {
      current: page,
      pageSize: pageSize,
      defaultPageSize: 50,
      total: logData.total,
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
      if (record.time === logData.list[index - 1].time && record.player === logData.list[index - 1].player) {
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
        <Card extra= {this.searchForm()}>
          <Row>
            <Col span={2}>
              <Card title="职业">
                {
                  playerPros.slice(1).map(value => {
                    return <p style={{textAlign: "center"}} key={value.name}>
                      <img src={require("../../assets/" + value.logo)} onClick={()=>this.queryClassTop(value.class)}/>
                    </p>
                  })
                }
              </Card>
            </Col>
            <Col span={7}>
              <Card title="伤害排行">
                <Table
                  bordered
                  size="small"
                  columns={this.columnsClassTop}
                  dataSource={classTop}
                  rowKey={(record) => {
                    return record.skill
                  }}
                  loading={loadingTop}
                  pagination={{
                    defaultPageSize: 50,
                    total: classTop.length,
                    pageSizeOptions: ['50', '100', '200', '500'],
                    showTotal: (total) => `共${total}条记录`,
                  }}
                />
              </Card>
            </Col>
            <Col span={15}>
              <Card title="原始日志">
                <Table
                  bordered
                  size="small"
                  columns={this.columns}
                  dataSource={logData.list}
                  rowKey={(record) => {
                    return record.id
                  }}
                  pagination={pagination}
                  loading={loading}
                  rowClassName={rowClassName}
                />
              </Card>
            </Col>
          </Row>
        </Card>
      </PageContainer>
    );
  }
}

export default Log;
