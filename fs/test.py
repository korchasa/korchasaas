import unittest
from tree import Tree
from var_dump import var_dump

class TestStringMethods(unittest.TestCase):

    def setUp(self):
        self.tree = Tree()

    def test_node_children(self):
        root = self.tree.node_children("/")
        self.assertIsInstance(root, dict)
        self.assertEqual(3, len(root))
        self.assertTrue('jobs_queue' in root)

        leaf = self.tree.node_children('/jobs_queue')
        self.assertIsInstance(leaf, dict)
        self.assertEqual(0, len(leaf))

if __name__ == '__main__':
    unittest.main()