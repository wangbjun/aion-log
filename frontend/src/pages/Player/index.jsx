import {Button, Card, Col, DatePicker, Form, Input, message, Row, Select, Statistic, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {Link} from 'umi';
import {playerPros} from "@/utils/utils";

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
        defaultSortOrder: 'ascend',
        sorter: function (a, b) {
          return a.name.localeCompare(b.name)
        },
        render: this.renderName
      },
      {
        title: "种族",
        dataIndex: 'type',
        key: 'type',
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
        dataIndex: 'pro',
        key: 'pro',
        sorter: function (a, b) {
          return a.pro - b.pro
        },
        render: function (value) {
          return <img src={require("../../assets/"+playerPros[value].logo)} width={35}/>
        }
      },
      {
        title: "最后更新时间",
        dataIndex: 'time',
        key: 'time',
        width: 180,
        sorter: function (a, b) {
          return moment(a.time).isAfter(moment(b.time))
        },
        render: function (value) {
          return moment(value).format("YYYY-MM-DD HH:mm:ss")
        }
      },
      {
        title: "操作",
        dataIndex: 'option',
        key: 'option',
        width: '30%',
        render: this.renderOption
      },
    ];
  }

  renderName = (value) => {
    return <Link to={`/log?player=${value}`}>{value}</Link>
  }

  renderOption = (value, row) => {
    return (<div>
      <a onClick={() => this.changeType(row, 1)}><Tag color="green">设为天族</Tag></a>
      <a onClick={() => this.changeType(row, 2)}><Tag color="blue">设为魔族</Tag></a>
      <a onClick={() => this.changeType(row, 0)}><Tag color="orange">设为其它</Tag></a>
    </div>)
  }

  async changeType(row, type) {
    const {dispatch} = this.props
    const result = await dispatch({
      type: 'global/changePlayerType',
      payload: {
        id: row.id,
        type: type
      }
    });
    if (result.code === 200) {
      message.success("操作成功")
      this.query()
    } else if (result.code === 405) {
      message.error("操作未授权")
    }
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
        pro: fieldValue.pro
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
        <Form.Item label="时间" name="time" style={{marginTop: "5px"}}>
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
          />
        </Form.Item>
        <Form.Item label="玩家" name="name" style={{marginTop: "5px"}}>
          <Input allowClear placeholder="请输入"/>
        </Form.Item>
        <Form.Item label="种族" name="type" style={{marginTop: "5px"}}>
          <Select
            allowClear
            showSearch
            style={{width: 150}}
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
        <Form.Item label="职业" name="pro" style={{marginTop: "5px"}}>
          <Select
            allowClear
            showSearch
            style={{width: 150}}
            placeholder="请选择种族"
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

  render() {
    const {playerList, loading} = this.props
    const statData = this.getStatData(playerList)
    return (
      <PageContainer>
        <Card extra={this.searchForm()}>
          <Row gutter={12}>
            <Col span={18}>
              <Table
                bordered
                size="small"
                columns={this.columns}
                dataSource={playerList}
                rowKey={(record) => {
                  return record.name + record.type
                }}
                pagination={{
                  defaultPageSize: 15,
                  hideOnSinglePage: true,
                  showTotal: (total) => `共${total}条记录`,
                }}
                loading={loading}
              />
            </Col>
            <Col span={6}>
              <Card title="概况">
                <Row gutter={24}>
                  <Col span={12}>
                    <Statistic title="天族" value={statData.tian} style={{padding: "24px"}}
                               valueStyle={{color: "green"}}/>
                  </Col>
                  <Col span={12}>
                    <Statistic title="总数" value={playerList.length} style={{padding: "24px"}}
                               valueStyle={{color: "red"}}/>
                  </Col>
                  <Col span={12}>
                    <Statistic title="魔族" value={statData.mo} style={{padding: "24px"}} valueStyle={{color: "blue"}}/>
                  </Col>
                  <Col span={12}>
                    <Statistic title="其它" value={statData.other} style={{padding: "24px"}}
                               valueStyle={{color: "orange"}}/>
                  </Col>
                </Row>
              </Card>
            </Col>
          </Row>
        </Card>
      </PageContainer>
    );
  }
}

export default Player;
