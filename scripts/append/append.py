import boto3
import numpy as np
from pprint import pprint
from argparse import ArgumentParser
import sys

log_client = boto3.client('logs')


def get_log_streams(lambda_id):
    group = '/aws/lambda/beldi-dev-{}'.format(lambda_id)
    r = log_client.describe_log_streams(logGroupName=group)
    r = r['logStreams']
    r = [x['logStreamName'] for x in r]
    return r


def delete_logs(lambda_id):
    group = '/aws/lambda/beldi-dev-{}'.format(lambda_id)
    try:
        log_client.delete_log_group(logGroupName=group)
    except:
        pass


def get_logs(lambda_id):
    group = '/aws/lambda/beldi-dev-{}'.format(lambda_id)
    streams = get_log_streams(lambda_id)
    res = []
    for stream in streams:
        r = log_client.get_log_events(logGroupName=group,
                                      logStreamName=stream)
        r = [e['message'].strip() for e in r['events']]
        res += r
    return res


## TODO: Modify the tags
def get_res(name):
    logs = get_logs(name)
    # print("\n\n\nLogs for:", name)
    # print('\n'.join(logs))
    tags = ["TPLRead", "TPLWrite", "Append", "Txn"]
    res = {}
    for tag in tags:
        res[tag] = []
    logs = list(filter(lambda x: 'DURATION' in x, logs))
    for log in logs:
        rs = log.strip().split()
        tag = rs[1]

        ## Parse the time and unit
        second_unit = rs[-1][-1:]
        assert(second_unit == "s")
        maybe_unit = rs[-1][-2]
        if (maybe_unit.isnumeric()):
            time = float(rs[-1][:-1])
            unit = second_unit
        else:
            time = float(rs[-1][:-2])
            unit = rs[-1][-2:]
        
        ## Adjust according to the unit
        if (unit == "s"):
            time *= 1000000
        elif (unit == "ms"):
            time *= 1000
        elif (unit == "\N{MICRO SIGN}s"):
            time *= 1
        else:
            print("Error: Unknown time unit:", unit)
            exit(1)
        res[tag].append(time)
    for k, v in res.items():
        v = np.array(v)
        res[k] = v[v < np.mean(v) * 2]  # Remove cold start for invocation
    p50p99 = {}
    for k, v in res.items():
        p50p99[k] = [np.percentile(v, 50), np.percentile(v, 99)]
    return p50p99


def main():
    parser = ArgumentParser()
    parser.add_argument("--command", required=True)
    args = parser.parse_args()
    if args.command == 'clean':
        delete_logs("bappend")
        delete_logs("append")
        delete_logs("tappend")
        return
    if args.command == 'run':
        # baseline = get_res("bappend")
        beldi = get_res("append")
        beldi_txn = get_res("tappend")
        # with open("result/append/append", "w") as f:
        f = sys.stdout
        f.write("#{:<19} {:<20} {:<20} {:<20} {:<20}\n".format("op",
                                                                            #  "Baseline", "Baseline 99",
                                                                                "Beldi", "Beldi 99",
                                                                                "Beldi-Txn", "Beldi-Txn 99"
                                                                                ))
        f.write("{:<20} {:<20} {:<20} {:<20} {:<20}\n".format("TPLRead",
                                                                            beldi["TPLRead"][0],
                                                                            beldi["TPLRead"][1],
                                                                            beldi_txn["TPLRead"][0],
                                                                            beldi_txn["TPLRead"][1]))
        f.write("{:<20} {:<20} {:<20} {:<20} {:<20}\n".format("TPLWrite",
                                                                            beldi["TPLWrite"][0],
                                                                            beldi["TPLWrite"][1],
                                                                            beldi_txn["TPLWrite"][0],
                                                                            beldi_txn["TPLWrite"][1]))
        f.write("{:<20} {:<20} {:<20} {:<20} {:<20}\n".format("Append",
                                                                            beldi["Append"][0],
                                                                            beldi["Append"][1],
                                                                            beldi_txn["Append"][0],
                                                                            beldi_txn["Append"][1]))
        f.write("{:<20} {:<20} {:<20} {:<20} {:<20}\n".format("Txn",
                                                                            beldi["Txn"][0],
                                                                            beldi["Txn"][1],
                                                                            beldi_txn["Txn"][0],
                                                                            beldi_txn["Txn"][1]))
        

if __name__ == "__main__":
    main()
