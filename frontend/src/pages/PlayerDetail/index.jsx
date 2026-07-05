import {Button, Card, DatePicker, Form, Input, Select, Table, Tag} from 'antd';
import React from 'react';
import {PageContainer} from '@ant-design/pro-layout';
import {connect} from "@/.umi/plugin-dva/exports";
import moment from "moment";
import {Link} from 'umi';
import {playerPros} from "@/utils/utils";
import "../../global.less"

const {RangePicker} = DatePicker
const {Option} = Select

@connect(
  state => ({
    ...state.global,
    loading: state.loading.effects["global/fetchPlayerList"]
  })
)
class PlayerDetail extends React.Component {
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
          return <img src={require("../../assets/" + playerPros[value].logo)} width={30} alt={playerPros[value].name}/>
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
          return (value * 100).toFixed(1) + "%"
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
    let st, et
    if (fieldValue.time) {
      st = fieldValue.time[0].format("YYYY-MM-DD HH:mm:ss")
      et = fieldValue.time[1].format("YYYY-MM-DD HH:mm:ss")
    }

    await dispatch({
      type: 'global/fetchPlayerList',
      payload: {
        st, et,
        name: fieldValue.name,
        type: fieldValue.type,
        class: fieldValue.class
      }
    });
  }

  onReset = async () => {
    this.formRef.current.resetFields();
    this.query().then()
  };

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
            ranges={{
              今天: [moment().startOf('day'), moment().endOf('day')],
              昨天: [moment().subtract(1, 'day').startOf('day'), moment().subtract(1, 'day').endOf('day')],
              前天: [moment().subtract(2, 'day').startOf('day'), moment().subtract(2, 'day').endOf('day')],
              最近3天: [moment().subtract(2, 'day').startOf('day'), moment().endOf('day')],
              最近7天: [moment().subtract(6, 'day').startOf('day'), moment().endOf('day')],
            }}
            allowClear
            showTime={{defaultValue: moment('00:00:00', 'HH:mm:ss')}}
            onChange={() => this.query()}
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
    return (
      <PageContainer title={false}>
        <Card className="battle-toolbar">
          {this.searchForm()}
        </Card>
        <Card title="玩家明细" className="battle-section-card">
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
              pageSizeOptions: ['20', '50', '100', '200', '500'],
              showTotal: (total) => `共${total}条记录`,
            }}
            loading={loading}
            scroll={{x: 900}}
          />
        </Card>
      </PageContainer>
    );
  }
}

export default PlayerDetail;
